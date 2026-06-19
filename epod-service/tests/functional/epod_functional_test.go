package functional

import (
    "net/http"
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestEPODFunctional_DatabaseAndAPI(t *testing.T) {
    resp, err := http.Get("http://localhost:8080/api/v1/epod/health")
    
    assert.NoError(t, err, "Functional Test Gagal: Koneksi ke server utama ditolak!")
    if err == nil {
        assert.Equal(t, http.StatusOK, resp.StatusCode)
    }
}