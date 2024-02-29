//go:build unit

package product_test

import "errors"

func (s *ProductSuite) TestAddMaintainer() {
	s.productClient.EXPECT().AddMaintainerToProduct(s.server, _productName, _userEmail).Return(nil).Once()
	s.renderer.EXPECT().RenderAddMaintainerToProduct(_productName, _userEmail).Times(1)

	err := s.handler.AddMaintainer(s.server.Name, _productName, _userEmail)
	s.NoError(err)
}

func (s *ProductSuite) TestAddMaintainer_ClientError() {
	expectedError := errors.New("client error")

	s.productClient.EXPECT().AddMaintainerToProduct(s.server, _productName, _userEmail).Return(expectedError).Once()

	err := s.handler.AddMaintainer(s.server.Name, _productName, _userEmail)
	s.Error(err, expectedError)
}
