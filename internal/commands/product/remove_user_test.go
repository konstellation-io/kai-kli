//go:build unit

package product_test

import "errors"

func (s *ProductSuite) TestRemoveUser() {
	s.productClient.EXPECT().RemoveUserFromProduct(s.server, _productName, _userEmail).Return(nil).Once()
	s.renderer.EXPECT().RenderRemoveUserFromProduct(_productName, _userEmail).Times(1)

	err := s.handler.RemoveUser(s.server.Name, _productName, _userEmail)
	s.NoError(err)
}

func (s *ProductSuite) TestRemoveUser_ClientError() {
	expectedError := errors.New("client error")

	s.productClient.EXPECT().RemoveUserFromProduct(s.server, _productName, _userEmail).Return(expectedError).Once()

	err := s.handler.RemoveUser(s.server.Name, _productName, _userEmail)
	s.Error(err, expectedError)
}
