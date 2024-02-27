//go:build unit

package product_test

import "errors"

func (s *ProductSuite) TestRemoveMaintainer() {
	s.productClient.EXPECT().RemoveMaintainerFromProduct(s.server, _productName, _userEmail).Return(nil).Once()
	s.renderer.EXPECT().RenderRemoveMaintainerFromProduct(_productName, _userEmail).Times(1)

	err := s.handler.RemoveMaintainer(s.server.Name, _productName, _userEmail)
	s.NoError(err)
}

func (s *ProductSuite) TestRemoveMaintainer_ClientError() {
	expectedError := errors.New("client error")

	s.productClient.EXPECT().RemoveMaintainerFromProduct(s.server, _productName, _userEmail).Return(expectedError).Once()

	err := s.handler.RemoveMaintainer(s.server.Name, _productName, _userEmail)
	s.Error(err, expectedError)
}
