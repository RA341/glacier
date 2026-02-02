package user

import v1 "github.com/ra341/glacier/generated/user/v1"

func (u *User) ToProto() *v1.User {
	return &v1.User{
		Id:       uint64(u.ID),
		Username: u.Username,
		Role:     u.Role.String(),
		//Password: u.EncryptedPassword, dont send the password
	}
}

func (u *User) FromProto(user *v1.User) error {
	roleString, err := RoleString(user.Role)
	if err != nil {
		return err
	}
	u.ID = uint(user.Id)
	u.Role = roleString
	u.Username = user.Username

	// this will be unencrypted
	u.EncryptedPassword = user.Password

	return nil
}
