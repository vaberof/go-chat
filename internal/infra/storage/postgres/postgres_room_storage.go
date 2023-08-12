package postgres

import (
	"errors"
	"github.com/vaberof/go-chat/internal/domain/room"
	"github.com/vaberof/go-chat/pkg/domain"
	"github.com/vaberof/go-chat/pkg/logging/logs"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"time"
)

type Room struct {
	Id        int64 `gorm:"primaryKey"`
	CreatorId int64
	Name      string
	Type      string
	Members   []*Member `gorm:"foreignKey:RoomId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Member struct {
	Id       int64 `gorm:"primaryKey"`
	UserId   int64
	RoomId   int64
	Nickname string
	Role     string
	JoinedAt time.Time
}

type roomStorageImpl struct {
	db     *gorm.DB
	logger *zap.SugaredLogger
}

func NewRoomStorage(db *gorm.DB, logs *logs.Logs) room.RoomStorage {
	loggerName := "room-storage"
	logger := logs.WithName(loggerName)

	return &roomStorageImpl{db: db, logger: logger}
}

func (storage *roomStorageImpl) Create(creatorId domain.UserId, name, roomType string, members []domain.UserId) (*room.Room, error) {
	postgresRoom := Room{
		CreatorId: int64(creatorId),
		Name:      name,
		Type:      roomType,
		Members:   make([]*Member, 0),
	}

	creatorUser, err := storage.getUser(creatorId)
	if err != nil {
		return nil, err
	}

	regularUsers, err := storage.getUsers(members)
	if err != nil {
		return nil, err
	}

	err = storage.db.Table("rooms").Create(&postgresRoom).Error
	if err != nil {
		storage.logger.Errorf("Failed to create a new room: %v", err)

		return nil, err
	}

	var postgresMembers []*Member

	adminMember := storage.convertUserToMember(creatorUser, &postgresRoom, room.AdminRole)
	regularMembers := storage.convertUsersToMembers(regularUsers, &postgresRoom, room.RegularRole)

	postgresMembers = append(postgresMembers, adminMember)
	postgresMembers = append(postgresMembers, regularMembers...)

	postgresRoom.Members = postgresMembers

	err = storage.db.Save(&postgresRoom).Error
	if err != nil {
		storage.logger.Errorf("Failed to save room in database: %v", err)

		return nil, err
	}

	return buildDomainRoom(&postgresRoom), nil
}

func (storage *roomStorageImpl) Get(roomId domain.RoomId) (*room.Room, error) {
	var postgresRoom Room

	err := storage.db.Table("rooms").Where("id = ?", roomId).First(&postgresRoom).Error
	if err != nil {
		storage.logger.Errorf("Failed to get a room: %v", err)

		return nil, err
	}

	return buildDomainRoom(&postgresRoom), nil
}

func (storage *roomStorageImpl) GetRooms(roomIds []domain.RoomId) ([]*room.Room, error) {
	var postgresRooms []*Room

	err := storage.db.Table("rooms").Where("id IN (?)", roomIds).Find(&postgresRooms).Error
	if err != nil {
		storage.logger.Errorf("Failed to get rooms: %v", err)

		return nil, err
	}

	return buildDomainRooms(postgresRooms), nil
}

// List returns user`s rooms according to given user id
func (storage *roomStorageImpl) List(userId domain.UserId) ([]*room.Room, error) {
	var user User

	err := storage.db.Preload("Rooms").Table("users").Where("id = ?", userId).First(&user).Error
	if err != nil {
		return nil, err
	}

	return buildDomainRooms(user.Rooms), nil
}

func (storage *roomStorageImpl) GetMembers(roomId domain.RoomId) ([]*room.Member, error) {
	var postgresMembers []*Member

	err := storage.db.Table("members").Where("room_id = ?", roomId).Find(&postgresMembers).Error
	if err != nil {
		storage.logger.Errorf("Failed to get members: %v", err)

		return nil, err
	}

	return buildDomainMembers(postgresMembers), nil
}

func (storage *roomStorageImpl) getUser(userId domain.UserId) (*User, error) {
	var postgresUser User

	err := storage.db.Preload("Rooms").Table("users").Where("id = ?", userId).First(&postgresUser).Error
	if err != nil {
		storage.logger.Errorf("Failed to get a user with id '%d': %v", userId, err)

		return nil, err
	}

	return &postgresUser, nil
}

func (storage *roomStorageImpl) getUsers(userIds []domain.UserId) ([]*User, error) {
	var postgresUsers []*User

	storage.db.Preload("Rooms").Table("users").Where("id IN(?)", userIds).Find(&postgresUsers)
	if len(postgresUsers) != len(userIds) {
		err := errors.New("users with given ids are not found")
		storage.logger.Errorf("Failed to get users: %v", err)

		return nil, err
	}

	return postgresUsers, nil
}

func (storage *roomStorageImpl) convertUsersToMembers(users []*User, room *Room, role string) []*Member {
	members := make([]*Member, len(users))

	for i := 0; i < len(members); i++ {
		members[i] = storage.convertUserToMember(users[i], room, role)
	}

	return members
}

func (storage *roomStorageImpl) convertUserToMember(user *User, room *Room, role string) *Member {
	return &Member{
		UserId:   int64(user.Id),
		Nickname: user.Username,
		Role:     role,
		JoinedAt: room.CreatedAt,
	}
}

func buildDomainRoom(postgresRoom *Room) *room.Room {
	return &room.Room{
		Id:        domain.RoomId(postgresRoom.Id),
		CreatorId: domain.UserId(postgresRoom.CreatorId),
		Name:      postgresRoom.Name,
		Type:      postgresRoom.Type,
		Members:   getMemberIds(postgresRoom.Members),
	}
}

func buildDomainRooms(postgresRooms []*Room) []*room.Room {
	domainRooms := make([]*room.Room, len(postgresRooms))

	for i := 0; i < len(domainRooms); i++ {
		domainRooms[i] = buildDomainRoom(postgresRooms[i])
	}

	return domainRooms
}

func buildDomainMember(postgresMember *Member) *room.Member {
	return &room.Member{
		UserId:   domain.UserId(postgresMember.UserId),
		RoomId:   domain.RoomId(postgresMember.RoomId),
		Nickname: postgresMember.Nickname,
		Role:     postgresMember.Role,
		JoinedAt: postgresMember.JoinedAt,
	}
}

func buildDomainMembers(postgresMembers []*Member) []*room.Member {
	domainMembers := make([]*room.Member, len(postgresMembers))

	for i := 0; i < len(domainMembers); i++ {
		domainMembers[i] = buildDomainMember(postgresMembers[i])
	}

	return domainMembers
}

func getMemberIds(members []*Member) []domain.UserId {
	roomIds := make([]domain.UserId, len(members))

	for i := 0; i < len(roomIds); i++ {
		roomIds[i] = domain.UserId(members[i].UserId)
	}

	return roomIds
}
