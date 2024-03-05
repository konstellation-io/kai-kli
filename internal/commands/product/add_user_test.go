//go:build unit

package product_test

import "errors"

const (
	_userEmail = "test@sample"
)

func (s *ProductSuite) TestAddUser() {
	s.productClient.EXPECT().AddUserToProduct(s.server, _productName, _userEmail).Return(nil).Once()
	s.renderer.EXPECT().RenderAddUserToProduct(_productName, _userEmail).Times(1)

	err := s.handler.AddUser(s.server.Name, _productName, _userEmail)
	s.NoError(err)
}

func (s *ProductSuite) TestAddUser_ClientError() {
	expectedError := errors.New("client error")

	s.productClient.EXPECT().AddUserToProduct(s.server, _productName, _userEmail).Return(expectedError).Once()

	err := s.handler.AddUser(s.server.Name, _productName, _userEmail)
	s.Error(err, expectedError)
}
