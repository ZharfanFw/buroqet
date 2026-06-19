package functional

import (
    "net/http"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestPricingFunctional_CalculateTariff(t *testing.T) {
    resp, err := http.Post("http://localhost:8080/pricing/calculate", "application/json", nil)
   
    assert.NoError(t, err, "Functional Test Gagal: Server utama belum aktif atau database belum siap!")
    if err == nil {
        assert.Equal(t, http.StatusOK, resp.StatusCode)
    }
}   