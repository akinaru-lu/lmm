package application

import (
	"context"

	"lmm/api/pkg/transaction"
	"lmm/api/service/user/application/command"
	"lmm/api/service/user/application/query"
	"lmm/api/service/user/domain"
	"lmm/api/service/user/domain/factory"
	"lmm/api/service/user/domain/model"
	"lmm/api/service/user/domain/repository"
	"lmm/api/service/user/domain/service"

	"github.com/pkg/errors"
)

// Service is a application service
type Service struct {
	encrypter          service.EncryptService
	factory            *factory.Factory
	transactionManager transaction.Manager
	userRepository     repository.UserRepository
}

// NewService creates a new Service pointer
func NewService(
	encrypter service.EncryptService,
	txManager transaction.Manager,
	userRepository repository.UserRepository,
) *Service {
	return &Service{
		encrypter:          encrypter,
		factory:            factory.NewFactory(encrypter, userRepository),
		transactionManager: txManager,
		userRepository:     userRepository,
	}
}

// RegisterNewUser registers new user
func (s *Service) RegisterNewUser(c context.Context, cmd command.Register) (int64, error) {
	var userID int64

	err := s.transactionManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		if user, err := s.userRepository.FindByName(tx, cmd.UserName); err != domain.ErrNoSuchUser {
			if user != nil {
				return domain.ErrUserNameAlreadyUsed
			}
			return errors.Wrap(err, "error occurred when checking user name duplication")
		}

		user, err := s.factory.NewUser(tx, cmd.UserName, cmd.EmailAddress, cmd.Password)
		if err != nil {
			return errors.Wrap(err, "invalid user")
		}

		if err := s.userRepository.Save(tx, user); err != nil {
			return errors.Wrap(err, "failed to save user")
		}

		userID = int64(user.ID())

		return nil
	}, nil)

	if err != nil {
		return 0, err
	}

	return userID, nil
}

// AssignRole handles command which operator assign user to role
func (s *Service) AssignRole(c context.Context, cmd command.AssignRole) error {
	return errors.New("not implemented")
	// return s.transactionManager.RunInTransaction(c, func(tx transaction.Transaction) error {
	// 	operator, err := s.userRepository.FindByName(tx, cmd.OperatorUser)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return errors.Wrap(domain.ErrNoSuchUser, err.Error())
	// 	}

	// 	user, err := s.userRepository.FindByName(tx, cmd.TargetUser)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return errors.Wrap(domain.ErrNoSuchUser, err.Error())
	// 	}

	// 	role := service.RoleAdapter(cmd.TargetRole)

	// 	return service.AssignUserRole(c, operator, user, role)
	// }, nil)
}

const maxCount uint = 100

func (s *Service) ViewAllUsersByOptions(c context.Context, query query.ViewAllUsers) ([]*model.UserDescriptor, uint, error) {
	panic("not implemented")
	// page, err := stringutil.ParseUint(query.Page)
	// if err != nil || page == 0 {
	// 	return nil, 0, errors.Wrap(domain.ErrInvalidPage, query.Page)
	// }

	// count, err := stringutil.ParseUint(query.Count)
	// if err != nil || count > maxCount {
	// 	return nil, 0, errors.Wrap(domain.ErrInvalidCount, query.Count)
	// }

	// order, err := s.mappingOrder(query.OrderBy, query.Order)
	// if err != nil {
	// 	return nil, 0, errors.Wrap(domain.ErrInvalidViewOrder, query.Order)
	// }

	// tx, err := s.transactionManager.Begin(c, nil)
	// if err != nil {
	// 	return nil, 0, err
	// }

	// users, i, err := s.userRepository.DescribeAll(tx, repository.DescribeAllOptions{
	// 	Page:  page,
	// 	Count: count,
	// 	Order: order,
	// })
	// if err != nil {
	// 	tx.Rollback()
	// 	return nil, 0, err
	// }

	// return users, i, nil
}

func (s *Service) mappingOrder(orderBy, order string) (repository.DescribeAllOrder, error) {
	switch orderBy + "_" + order {
	case "name_asc":
		return repository.DescribeAllOrderByNameAsc, nil
	case "name_desc":
		return repository.DescribeAllOrderByNameDesc, nil
	case "registered_date_asc":
		return repository.DescribeAllOrderByRegisteredDateAsc, nil
	case "registered_date_desc":
		return repository.DescribeAllOrderByRegisteredDateDesc, nil
	case "role_asc":
		return repository.DescribeAllOrderByRoleAsc, nil
	case "role_desc":
		return repository.DescribeAllOrderByRoleDesc, nil
	default:
		return repository.DescribeAllOrder(-1), domain.ErrInvalidViewOrder
	}
}

// UserChangePassword supports a application to chagne user's password
func (s *Service) UserChangePassword(c context.Context, cmd command.ChangePassword) error {
	hashedPassword, err := s.factory.NewPassword(cmd.NewPassword)
	if err != nil {
		return errors.Wrap(err, "invalid password")
	}

	return s.transactionManager.RunInTransaction(c, func(tx transaction.Transaction) error {
		user, err := s.userRepository.FindByName(tx, cmd.User)
		if err != nil {
			return errors.Wrap(domain.ErrNoSuchUser, err.Error())
		}

		if !s.encrypter.Verify(cmd.OldPassword, user.Password()) {
			return domain.ErrUserPassword
		}

		if err := user.ChangePassword(hashedPassword); err != nil {
			return errors.Wrap(err, "failed to change password")
		}

		if err := user.ChangeToken(s.factory.NewToken()); err != nil {
			return errors.Wrap(err, "failed to change token")
		}

		if err := s.userRepository.Save(tx, user); err != nil {
			return errors.Wrap(err, "failed to save user after password and token changed")
		}

		return nil
	}, nil)
}
