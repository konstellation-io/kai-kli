//go:build unit

package productconfiguration_test

import productconfiguration "github.com/konstellation-io/kli/internal/services/product_configuration"

func (ch *KaiProductConfigurationTest) TestGetProductVersion_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	productVersion := kaiProductConfig.GetProductVersion()

	// THEN
	ch.Require().Equal(_defaultProductVersion, productVersion)
}

func (ch *KaiProductConfigurationTest) TestUpdateProductVersion_ExpectOk() {
	// GIVEN
	newVersion := "v1.2.3"
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	err := kaiProductConfig.UpdateProductVersion(newVersion)

	// THEN
	ch.Require().NoError(err)
	ch.Assert().Equal(newVersion, kaiProductConfig.GetProductVersion())
}

func (ch *KaiProductConfigurationTest) TestUpdateProductVersionAndDescription_ExpectOk() {
	// GIVEN
	newVersion := "v1.2.3"
	newDescription := "New description"
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	err := kaiProductConfig.UpdateProductVersion(newVersion, newDescription)

	// THEN
	ch.Require().NoError(err)
	ch.Require().Equal(newVersion, kaiProductConfig.GetProductVersion())
	ch.Require().Equal(newDescription, kaiProductConfig.GetProductDescription())
}

func (ch *KaiProductConfigurationTest) TestUpdateProductVersion_InvalidVersion_ExpectError() {
	// GIVEN
	newVersion := "v.1.2.3"
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	err := kaiProductConfig.UpdateProductVersion(newVersion)

	// THEN
	ch.Assert().Error(err)
	ch.Assert().ErrorIs(err, productconfiguration.ErrInvalidVersion)
}

func (ch *KaiProductConfigurationTest) TestGetProductDescription_ExpectOk() {
	// GIVEN
	kaiProductConfig := ch.getDefaultConfiguration()

	// WHEN
	productDescription := kaiProductConfig.GetProductDescription()

	// THEN
	ch.Require().Equal(_defaultProductDescription, productDescription)
}
