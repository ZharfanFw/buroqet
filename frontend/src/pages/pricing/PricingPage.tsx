import { useState } from "react";

// 🌟 Interface untuk mendefinisikan tipe data agar TypeScript tidak protes
interface RateData {
  serviceName: string;
  serviceCode: string;
  price: number;
  estimation: string;
}

export default function PricingPage() {
  // Nilai awal diarahkan ke data yang pasti ada harganya di Supabase
  const [origin, setOrigin] = useState("10110");
  const [destination, setDestination] = useState("40115");
  const [weight, setWeight] = useState<number | string>(1);
  const [dimensions, setDimensions] = useState("");

  const [loading, setLoading] = useState(false);
  const [rates, setRates] = useState<RateData[] | null>(null);

  const handleCalculate = async () => {
    setLoading(true);
    setRates(null);

    try {
      // Bedah string dimensi "10x15x20"
      const dims = dimensions.split("x").map((d) => parseFloat(d.trim()));
      const length_cm = dims[0] || 10;
      const width_cm = dims[1] || 10;
      const height_cm = dims[2] || 10;

      // Tembak API Golang
      const response = await fetch("/api/pricing/calculate", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          origin_postal_code: origin,
          destination_postal_code: destination,
          service_type: "REG",
          weight_kg: parseFloat(weight.toString()),
          length_cm: length_cm,
          width_cm: width_cm,
          height_cm: height_cm,
          use_insurance: false,
        }),
      });

      if (!response.ok) throw new Error("Gagal mengambil data ongkir");

      const data = await response.json();

      console.log("ISI DATA DARI GOLANG:", data);

      // Simpan data dari Backend ke state UI
      setRates([
        {
          serviceName: "Reguler (REG)",
          serviceCode: "REG",
          // 👇 UBAH DUA BARIS INI (tambahkan .data di tengahnya)
          price: data.data?.total || 0,
          estimation: data.data?.estimated || "2-3 Hari",
        },
      ]);
    } catch (error) {
      console.error(error);
      alert("Terjadi kesalahan saat menghitung ongkir.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ display: "flex", flexDirection: "column", gap: "20px" }}>
      <div className="page-header" style={{ marginBottom: 0 }}>
        <h1>💰 Cek Ongkir (Pricing)</h1>
        <p>Kalkulasi estimasi biaya pengiriman berdasarkan berat dan dimensi</p>
      </div>

      <div
        style={{ display: "grid", gridTemplateColumns: "1fr 1fr", gap: "20px" }}
      >
        {/* Calculator Form */}
        <div
          className="card"
          style={{ display: "flex", flexDirection: "column", gap: "16px" }}
        >
          <h3
            style={{ fontSize: "16px", fontWeight: 600, marginBottom: "4px" }}
          >
            Detail Pengiriman
          </h3>

          <div
            style={{
              display: "grid",
              gridTemplateColumns: "1fr 1fr",
              gap: "12px",
            }}
          >
            <div>
              <label
                style={{
                  display: "block",
                  fontSize: "12px",
                  fontWeight: 600,
                  color: "var(--text-secondary)",
                  marginBottom: "6px",
                }}
              >
                Kota Asal
              </label>
              <select
                className="input"
                value={origin}
                onChange={(e) => setOrigin(e.target.value)}
              >
                <option value="10110">Jakarta (10110)</option>
                <option value="40115">Bandung (40115)</option>
                <option value="60000">Surabaya (60000)</option>
              </select>
            </div>
            <div>
              <label
                style={{
                  display: "block",
                  fontSize: "12px",
                  fontWeight: 600,
                  color: "var(--text-secondary)",
                  marginBottom: "6px",
                }}
              >
                Kota Tujuan
              </label>
              <select
                className="input"
                value={destination}
                onChange={(e) => setDestination(e.target.value)}
              >
                <option value="40115">Bandung (40115)</option>
                <option value="60000">Surabaya (60000)</option>
                <option value="10110">Jakarta (10110)</option>
              </select>
            </div>
          </div>

          <div
            style={{
              display: "grid",
              gridTemplateColumns: "1fr 1fr",
              gap: "12px",
            }}
          >
            <div>
              <label
                style={{
                  display: "block",
                  fontSize: "12px",
                  fontWeight: 600,
                  color: "var(--text-secondary)",
                  marginBottom: "6px",
                }}
              >
                Berat (kg)
              </label>
              <input
                type="number"
                className="input"
                value={weight}
                onChange={(e) => setWeight(e.target.value)}
                min={1}
              />
            </div>
            <div>
              <label
                style={{
                  display: "block",
                  fontSize: "12px",
                  fontWeight: 600,
                  color: "var(--text-secondary)",
                  marginBottom: "6px",
                }}
              >
                Dimensi (PxLxT) cm
              </label>
              <input
                type="text"
                className="input"
                placeholder="Contoh: 10x10x10"
                value={dimensions}
                onChange={(e) => setDimensions(e.target.value)}
              />
            </div>
          </div>

          <button
            className="btn btn-primary"
            style={{
              width: "100%",
              justifyContent: "center",
              marginTop: "10px",
            }}
            onClick={handleCalculate}
            disabled={loading}
          >
            {loading ? "Menghitung..." : "Hitung Ongkos Kirim"}
          </button>
        </div>

        {/* Results */}
        <div
          className="card"
          style={{
            display: "flex",
            flexDirection: "column",
            background: rates ? "var(--bg-surface)" : "var(--bg-hover)",
          }}
        >
          <h3
            style={{ fontSize: "16px", fontWeight: 600, marginBottom: "16px" }}
          >
            Hasil Kalkulasi
          </h3>

          {!rates && !loading && (
            <div
              className="empty-state"
              style={{ padding: "40px 20px", margin: "auto" }}
            >
              <div className="empty-icon" style={{ opacity: 0.3 }}>
                🧾
              </div>
              <p>Masukkan detail pengiriman untuk melihat estimasi harga.</p>
            </div>
          )}

          {loading && (
            <div
              style={{
                display: "flex",
                justifyContent: "center",
                alignItems: "center",
                flex: 1,
              }}
            >
              <div
                className="spinner"
                style={{
                  borderColor: "var(--border)",
                  borderTopColor: "var(--primary)",
                }}
              ></div>
            </div>
          )}

          {rates && !loading && (
            <div
              style={{ display: "flex", flexDirection: "column", gap: "12px" }}
            >
              {rates.map((rate, index) => (
                <div
                  key={index}
                  style={{
                    border:
                      index === 0
                        ? "1px solid var(--primary)"
                        : "1px solid var(--border)",
                    borderRadius: "12px",
                    padding: "16px",
                    background:
                      index === 0
                        ? "rgba(121, 174, 111, 0.05)"
                        : "var(--bg-surface)",
                  }}
                >
                  <div
                    style={{
                      display: "flex",
                      justifyContent: "space-between",
                      alignItems: "center",
                      marginBottom: "8px",
                    }}
                  >
                    <div
                      style={{
                        fontWeight: 600,
                        color: index === 0 ? "var(--primary-dark)" : "inherit",
                      }}
                    >
                      {rate.serviceName}
                    </div>
                    <div style={{ fontSize: "18px", fontWeight: 700 }}>
                      Rp {rate.price.toLocaleString("id-ID")}
                    </div>
                  </div>
                  <div
                    style={{ fontSize: "12px", color: "var(--text-secondary)" }}
                  >
                    Estimasi tiba: {rate.estimation}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      </div>
    </div>
  );
}
