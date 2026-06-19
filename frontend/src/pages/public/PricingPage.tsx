import { useState } from 'react';

export default function PricingPage() {
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState(false);
  const [selectedService, setSelectedService] = useState<string | null>(null);

  const handleCalculate = () => {
    setLoading(true);
    // Simulasi loading 800ms
    setTimeout(() => {
      setLoading(false);
      setResult(true);
    }, 800);
  };

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col items-center py-12 px-4 sm:px-6 lg:px-8 font-sans">
      
      {/* Header Container - Centered */}
      <div className="w-full max-w-3xl text-center mb-10">
        <h1 className="text-4xl md:text-5xl font-extrabold text-slate-900 mb-4 tracking-tight">
          💰 Cek Ongkir
        </h1>
        <p className="text-lg text-slate-600">
          Kalkulasi estimasi biaya pengiriman berdasarkan berat dan dimensi
        </p>
      </div>

      {/* Main Content Container - Centered Card */}
      <div className="w-full max-w-4xl bg-white rounded-2xl shadow-sm border-2 border-slate-200 p-6 md:p-8">
        <div className="grid grid-cols-1 md:grid-cols-2 gap-8 md:gap-12">
          
          {/* Calculator Form */}
          <div className="flex flex-col gap-5">
            <h3 className="text-xl font-bold text-slate-800 border-b-2 border-slate-100 pb-2">Detail Pengiriman</h3>
            
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-semibold text-slate-600 mb-2">Kota Asal</label>
                <select className="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 text-slate-800 rounded-xl focus:bg-white focus:outline-none focus:border-[#79ae6f] focus:ring-4 focus:ring-[#79ae6f]/20 transition-all font-medium appearance-none">
                  <option>Jakarta</option>
                  <option>Bandung</option>
                  <option>Surabaya</option>
                </select>
              </div>
              <div>
                <label className="block text-sm font-semibold text-slate-600 mb-2">Kota Tujuan</label>
                <select className="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 text-slate-800 rounded-xl focus:bg-white focus:outline-none focus:border-[#79ae6f] focus:ring-4 focus:ring-[#79ae6f]/20 transition-all font-medium appearance-none">
                  <option>Surabaya</option>
                  <option>Bandung</option>
                  <option>Jakarta</option>
                </select>
              </div>
            </div>

            <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label className="block text-sm font-semibold text-slate-600 mb-2">Berat (kg)</label>
                <input 
                  type="number" 
                  className="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 text-slate-800 rounded-xl focus:bg-white focus:outline-none focus:border-[#79ae6f] focus:ring-4 focus:ring-[#79ae6f]/20 transition-all font-medium" 
                  defaultValue={1} 
                  min={1} 
                />
              </div>
              <div>
                <label className="block text-sm font-semibold text-slate-600 mb-2">Dimensi (PxLxT) cm</label>
                <input 
                  type="text" 
                  className="w-full px-4 py-3 bg-slate-50 border-2 border-slate-200 text-slate-800 rounded-xl focus:bg-white focus:outline-none focus:border-[#79ae6f] focus:ring-4 focus:ring-[#79ae6f]/20 transition-all font-medium placeholder-slate-400" 
                  placeholder="Contoh: 10x10x10" 
                />
              </div>
            </div>

            <button 
              className="mt-2 w-full bg-[#64965a] hover:bg-[#53804a] text-white px-8 py-3.5 rounded-xl font-semibold transition-colors disabled:opacity-60 disabled:cursor-not-allowed flex items-center justify-center shadow-sm hover:shadow-md"
              onClick={handleCalculate}
              disabled={loading}
            >
              {loading ? (
                <svg className="animate-spin h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              ) : 'Hitung Ongkos Kirim'}
            </button>
          </div>

          {/* Results Section */}
          <div className={`flex flex-col rounded-2xl p-6 transition-colors duration-300 ${result ? 'bg-slate-50' : 'bg-slate-50 border-2 border-dashed border-slate-200'}`}>
            <h3 className="text-xl font-bold text-slate-800 border-b-2 border-slate-100 pb-2 mb-4">Hasil Kalkulasi</h3>
            
            {/* Empty State */}
            {!result && !loading && (
              <div className="flex flex-col items-center justify-center flex-1 text-center opacity-60 min-h-[200px]">
                <span className="text-5xl mb-4">🧾</span>
                <p className="text-slate-600 font-medium">Masukkan detail pengiriman untuk melihat estimasi harga.</p>
              </div>
            )}

            {/* Loading Spinner */}
            {loading && (
              <div className="flex items-center justify-center flex-1 min-h-[200px]">
                <svg className="animate-spin h-10 w-10 text-[#64965a]" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
            )}

            {/* Results Data */}
            {result && !loading && (
              <div className="flex flex-col gap-4">
                {/* Primary Result (Highlighted) */}
                <div 
                  onClick={() => setSelectedService('REG')}
                  className={`cursor-pointer border-2 rounded-xl p-4 relative overflow-hidden transition-all ${selectedService === 'REG' ? 'border-[#64965a] bg-[#64965a]/10 shadow-md ring-2 ring-[#64965a]/20 scale-[1.02]' : 'border-[#e4ece3] bg-white shadow-sm hover:border-[#79ae6f] hover:bg-slate-50'}`}
                >
                  <div className="absolute top-0 right-0 bg-[#64965a] text-white text-[10px] font-bold px-3 py-1 rounded-bl-lg shadow-sm">TERCEPAT</div>
                  <div className="flex justify-between items-center mb-1 mt-1">
                    <div className={`font-bold ${selectedService === 'REG' ? 'text-[#53804a]' : 'text-slate-700'}`}>Reguler (REG)</div>
                    <div className="text-xl font-extrabold text-slate-800">Rp 15.000</div>
                  </div>
                  <div className={`text-sm font-medium ${selectedService === 'REG' ? 'text-[#64965a]' : 'text-slate-500'}`}>Estimasi tiba: 2-3 Hari</div>
                </div>
                
                {/* Secondary Result 1 */}
                <div 
                  onClick={() => setSelectedService('YES')}
                  className={`cursor-pointer border-2 rounded-xl p-4 relative overflow-hidden transition-all ${selectedService === 'YES' ? 'border-[#64965a] bg-[#64965a]/10 shadow-md ring-2 ring-[#64965a]/20 scale-[1.02]' : 'border-[#e4ece3] bg-white shadow-sm hover:border-[#79ae6f] hover:bg-slate-50'}`}
                >
                  <div className="flex justify-between items-center mb-1">
                    <div className={`font-bold ${selectedService === 'YES' ? 'text-[#53804a]' : 'text-slate-700'}`}>Next Day (YES)</div>
                    <div className="text-xl font-extrabold text-slate-800">Rp 22.000</div>
                  </div>
                  <div className={`text-sm font-medium ${selectedService === 'YES' ? 'text-[#64965a]' : 'text-slate-500'}`}>Estimasi tiba: 1 Hari (Besok)</div>
                </div>
                
                {/* Secondary Result 2 */}
                <div 
                  onClick={() => setSelectedService('CARGO')}
                  className={`cursor-pointer border-2 rounded-xl p-4 relative overflow-hidden transition-all ${selectedService === 'CARGO' ? 'border-[#64965a] bg-[#64965a]/10 shadow-md ring-2 ring-[#64965a]/20 scale-[1.02]' : 'border-[#e4ece3] bg-white shadow-sm hover:border-[#79ae6f] hover:bg-slate-50'}`}
                >
                  <div className="flex justify-between items-center mb-1">
                    <div className={`font-bold ${selectedService === 'CARGO' ? 'text-[#53804a]' : 'text-slate-700'}`}>Kargo (CARGO)</div>
                    <div className="text-xl font-extrabold text-slate-800">Rp 45.000</div>
                  </div>
                  <div className={`text-sm font-medium ${selectedService === 'CARGO' ? 'text-[#64965a]' : 'text-slate-500'}`}>Estimasi tiba: 4-7 Hari • Min 10kg</div>
                </div>

                {selectedService && (
                  <button className="mt-4 w-full bg-[#64965a] hover:bg-[#53804a] text-white px-6 py-3.5 rounded-xl font-bold shadow-md hover:shadow-lg transform hover:-translate-y-0.5 transition-all text-sm">
                    Pilih Layanan Ini
                  </button>
                )}
              </div>
            )}
          </div>

        </div>
      </div>
    </div>
  );
}