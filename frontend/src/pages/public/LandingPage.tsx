import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';

export default function LandingPage() {
  const [trackId, setTrackId] = useState('');
  const navigate = useNavigate();

  const handleTrack = (e: React.FormEvent) => {
    e.preventDefault();
    if (trackId) {
      navigate(`/tracking?awb=${trackId}`);
    }
  };

  return (
    <div className="bg-[#f7faf7] min-h-screen">
      {/* Hero Section */}
      <section className="pt-32 pb-20 px-4 sm:px-6 lg:px-8 max-w-7xl mx-auto relative">
        <div className="absolute top-20 right-0 w-[500px] h-[500px] bg-[#64965a]/10 rounded-full blur-3xl -z-10"></div>
        <div className="absolute bottom-0 left-0 w-[400px] h-[400px] bg-[#79ae6f]/10 rounded-full blur-3xl -z-10"></div>
        
        <div className="flex flex-col lg:flex-row items-center gap-16">
          <div className="flex-1 text-center lg:text-left">
            <div className="inline-block px-4 py-1.5 bg-green-50 text-green-600 rounded-full text-xs font-bold uppercase tracking-wider mb-6 border border-green-200">
              #1 Tech-Driven Logistics
            </div>
            <h1 className="text-5xl lg:text-7xl font-extrabold text-slate-900 leading-tight mb-6 tracking-tight">
              Kirim Paket <br />
              <span className="text-[#64965a]">Lebih Cepat</span> & <br />
              <span className="text-[#64965a]">Lebih Aman.</span>
            </h1>
            <p className="text-lg md:text-xl text-slate-500 mb-10 max-w-2xl mx-auto lg:mx-0 font-medium">
              Buroqet adalah partner logistik masa depan Anda. Lacak paket secara real-time dan nikmati pengalaman pengiriman tanpa hambatan dengan teknologi microservices.
            </p>
            
            <div className="flex flex-col sm:flex-row items-center gap-4 justify-center lg:justify-start">
              <Link to="/customer/register" className="w-full sm:w-auto px-8 py-4 bg-[#64965a] hover:bg-[#53804a] text-white font-extrabold rounded-2xl shadow-xl hover:shadow-2xl transform hover:-translate-y-1 transition-all text-center">
                Mulai Kirim Paket
              </Link>
              <a href="#tracking-box" className="w-full sm:w-auto px-8 py-4 bg-white border border-[#e4ece3] hover:bg-slate-50 text-slate-700 font-extrabold rounded-2xl shadow-sm hover:shadow-md transform hover:-translate-y-1 transition-all text-center">
                Lacak Resi
              </a>
            </div>
          </div>

          <div className="flex-1 relative w-full max-w-lg lg:max-w-none">
            {/* Elegant Logistics Illustration Placeholder (CSS Art/Box) */}
            <div className="relative bg-white p-6 rounded-[2.5rem] shadow-2xl border border-[#e4ece3] rotate-2 hover:rotate-0 transition-transform duration-500">
              <div className="absolute -top-6 -right-6 w-24 h-24 bg-yellow-400 rounded-full blur-xl opacity-50"></div>
              <img src="https://images.unsplash.com/photo-1566576721346-d4a3b4eaeb55?q=80&w=1000&auto=format&fit=crop" alt="Logistics Box" className="rounded-3xl w-full h-80 object-cover" />
              <div className="absolute -bottom-6 -left-6 bg-white p-4 rounded-2xl shadow-xl border border-[#e4ece3] flex items-center gap-4 animate-bounce">
                <div className="w-12 h-12 bg-green-50 text-green-500 rounded-full flex items-center justify-center text-xl">✅</div>
                <div>
                  <p className="text-xs text-slate-500 font-bold uppercase">Status</p>
                  <p className="text-sm font-extrabold text-slate-800">Paket Diterima</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      {/* Tracking Box */}
      <section id="tracking-box" className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 relative -mt-10 z-20">
        <div className="bg-white rounded-[2rem] p-8 md:p-10 shadow-2xl border border-[#e4ece3]">
          <h2 className="text-2xl font-extrabold text-slate-800 mb-6 text-center">Lacak Pengiriman Anda</h2>
          <form className="flex flex-col sm:flex-row gap-4" onSubmit={handleTrack}>
            <div className="relative flex-1">
              <span className="absolute inset-y-0 left-0 flex items-center pl-5 text-slate-400 text-xl">📍</span>
              <input 
                type="text" 
                placeholder="Masukkan No. Resi (Contoh: BQ-2024-JKT-102)" 
                className="w-full pl-14 pr-5 py-4 rounded-2xl border border-[#e4ece3] bg-slate-50 text-slate-800 font-bold placeholder:text-slate-400 focus:outline-none focus:ring-4 focus:ring-[#64965a]/20 focus:border-[#64965a] transition-all"
                value={trackId}
                onChange={(e) => setTrackId(e.target.value)}
              />
            </div>
            <button 
              type="submit"
              className="px-10 py-4 rounded-2xl bg-[#64965a] text-white font-extrabold hover:bg-[#53804a] transition-all shadow-lg hover:shadow-xl transform hover:-translate-y-1"
            >
              Lacak Paket
            </button>
          </form>
        </div>
      </section>

      {/* Statistics */}
      <section className="py-20 px-4">
        <div className="max-w-7xl mx-auto grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
          <div>
            <h3 className="text-4xl lg:text-5xl font-extrabold text-[#64965a] mb-2">99.9%</h3>
            <p className="text-slate-500 font-bold uppercase tracking-wider text-sm">On-time Delivery</p>
          </div>
          <div>
            <h3 className="text-4xl lg:text-5xl font-extrabold text-[#64965a] mb-2">5M+</h3>
            <p className="text-slate-500 font-bold uppercase tracking-wider text-sm">Paket Terkirim</p>
          </div>
          <div>
            <h3 className="text-4xl lg:text-5xl font-extrabold text-[#64965a] mb-2">12K+</h3>
            <p className="text-slate-500 font-bold uppercase tracking-wider text-sm">Mitra Kurir</p>
          </div>
          <div>
            <h3 className="text-4xl lg:text-5xl font-extrabold text-[#64965a] mb-2">24/7</h3>
            <p className="text-slate-500 font-bold uppercase tracking-wider text-sm">Customer Support</p>
          </div>
        </div>
      </section>

      {/* Why Choose Us */}
      <section id="about" className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p className="text-[#64965a] font-bold tracking-widest uppercase mb-2">Keunggulan Buroqet</p>
          <h2 className="text-3xl md:text-4xl font-extrabold text-slate-800 mb-16">Mengapa Memilih Kami?</h2>
          
          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            <div className="bg-[#f7faf7] p-10 rounded-[2rem] border border-[#e4ece3] hover:-translate-y-2 transition-transform duration-300">
              <div className="w-16 h-16 bg-white text-[#64965a] rounded-2xl flex items-center justify-center text-3xl mb-6 shadow-sm mx-auto">
                🛰️
              </div>
              <h3 className="text-xl font-bold text-slate-800 mb-4">Real-time Tracking</h3>
              <p className="text-slate-500 font-medium leading-relaxed">Pantau posisi pasti paket Anda setiap saat tanpa delay berkat teknologi event-streaming mutakhir dari Buroqet.</p>
            </div>
            
            <div className="bg-[#f7faf7] p-10 rounded-[2rem] border border-[#e4ece3] hover:-translate-y-2 transition-transform duration-300">
              <div className="w-16 h-16 bg-white text-yellow-500 rounded-2xl flex items-center justify-center text-3xl mb-6 shadow-sm mx-auto">
                ⚡
              </div>
              <h3 className="text-xl font-bold text-slate-800 mb-4">Pengiriman Ekstra Cepat</h3>
              <p className="text-slate-500 font-medium leading-relaxed">Sistem pintar kami mengalokasikan kurir terdekat untuk menjemput paket Anda dalam hitungan menit.</p>
            </div>

            <div className="bg-[#f7faf7] p-10 rounded-[2rem] border border-[#e4ece3] hover:-translate-y-2 transition-transform duration-300">
              <div className="w-16 h-16 bg-white text-blue-500 rounded-2xl flex items-center justify-center text-3xl mb-6 shadow-sm mx-auto">
                🛡️
              </div>
              <h3 className="text-xl font-bold text-slate-800 mb-4">Aman & Terpercaya</h3>
              <p className="text-slate-500 font-medium leading-relaxed">Dilengkapi dengan fitur ePOD (Electronic Proof of Delivery) beserta foto bukti serah terima paket yang sah.</p>
            </div>
          </div>
        </div>
      </section>

      {/* Services */}
      <section id="services" className="py-20 bg-[#f7faf7]">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="text-center mb-16">
            <p className="text-[#64965a] font-bold tracking-widest uppercase mb-2">Layanan Kami</p>
            <h2 className="text-3xl md:text-4xl font-extrabold text-slate-800 mb-4">Solusi Pengiriman Terbaik</h2>
          </div>

          <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
            {/* REGULAR */}
            <div className="bg-white rounded-[2rem] p-10 border border-[#e4ece3] hover:border-[#64965a] transition-all duration-300 shadow-sm hover:shadow-xl group">
              <div className="w-14 h-14 bg-slate-50 text-slate-700 rounded-2xl flex items-center justify-center text-2xl mb-8 group-hover:bg-[#64965a] group-hover:text-white transition-colors">📦</div>
              <h3 className="text-2xl font-extrabold text-slate-800 mb-3">Regular (REG)</h3>
              <p className="text-slate-500 font-medium mb-8 min-h-[60px]">Layanan standar yang terjangkau dengan jangkauan ke seluruh pelosok Indonesia.</p>
              <ul className="space-y-3 mb-8">
                <li className="flex items-center text-slate-600 font-medium"><span className="text-[#64965a] mr-3 font-bold">✓</span> Estimasi 2-4 Hari</li>
                <li className="flex items-center text-slate-600 font-medium"><span className="text-[#64965a] mr-3 font-bold">✓</span> Tracking Real-time</li>
              </ul>
              <div className="text-2xl font-extrabold text-[#64965a]">Rp 9.000<span className="text-sm text-slate-400 font-medium">/kg</span></div>
            </div>

            {/* EXPRESS / INSTANT */}
            <div className="bg-[#64965a] rounded-[2rem] p-10 border border-[#64965a] text-white shadow-2xl relative transform md:-translate-y-4">
              <div className="absolute top-6 right-6 bg-yellow-400 text-yellow-900 text-xs font-extrabold px-3 py-1.5 rounded-lg">TERPOPULER</div>
              <div className="w-14 h-14 bg-white/20 rounded-2xl flex items-center justify-center text-2xl mb-8">⚡</div>
              <h3 className="text-2xl font-extrabold text-white mb-3">Express (EXP)</h3>
              <p className="text-white/80 font-medium mb-8 min-h-[60px]">Pengiriman kilat untuk paket mendesak. Sampai di hari yang sama atau besok.</p>
              <ul className="space-y-3 mb-8">
                <li className="flex items-center text-white font-medium"><span className="text-yellow-300 mr-3 font-bold">✓</span> Estimasi 1 Hari / Same-day</li>
                <li className="flex items-center text-white font-medium"><span className="text-yellow-300 mr-3 font-bold">✓</span> Prioritas Penjemputan</li>
              </ul>
              <div className="text-2xl font-extrabold text-yellow-300">Rp 18.000<span className="text-sm text-white/70 font-medium">/kg</span></div>
            </div>

            {/* CARGO */}
            <div className="bg-white rounded-[2rem] p-10 border border-[#e4ece3] hover:border-[#64965a] transition-all duration-300 shadow-sm hover:shadow-xl group">
              <div className="w-14 h-14 bg-slate-50 text-slate-700 rounded-2xl flex items-center justify-center text-2xl mb-8 group-hover:bg-[#64965a] group-hover:text-white transition-colors">🚛</div>
              <h3 className="text-2xl font-extrabold text-slate-800 mb-3">Cargo (CARGO)</h3>
              <p className="text-slate-500 font-medium mb-8 min-h-[60px]">Solusi hemat untuk pengiriman barang besar atau dalam jumlah banyak (&gt;10kg).</p>
              <ul className="space-y-3 mb-8">
                <li className="flex items-center text-slate-600 font-medium"><span className="text-[#64965a] mr-3 font-bold">✓</span> Estimasi 4-7 Hari</li>
                <li className="flex items-center text-slate-600 font-medium"><span className="text-[#64965a] mr-3 font-bold">✓</span> Tarif Ekonomis</li>
              </ul>
              <div className="text-2xl font-extrabold text-[#64965a]">Rp 25.000<span className="text-sm text-slate-400 font-medium">/10kg</span></div>
            </div>
          </div>
        </div>
      </section>

      {/* Shipping Flow */}
      <section className="py-20 bg-white">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
          <p className="text-[#64965a] font-bold tracking-widest uppercase mb-2">Cara Kerja</p>
          <h2 className="text-3xl md:text-4xl font-extrabold text-slate-800 mb-16">Alur Pengiriman Buroqet</h2>
          <div className="flex flex-col md:flex-row justify-center items-center gap-8 md:gap-12 relative">
            {/* Connecting Line */}
            <div className="hidden md:block absolute top-1/2 left-20 right-20 h-0.5 bg-[#e4ece3] -z-10"></div>
            
            <div className="bg-white p-6 rounded-3xl w-64 border border-[#e4ece3] shadow-sm relative z-10">
              <div className="w-16 h-16 bg-[#64965a] text-white rounded-full flex items-center justify-center text-2xl font-bold mx-auto mb-4 border-4 border-white shadow-md">1</div>
              <h4 className="text-lg font-bold text-slate-800 mb-2">Buat Order</h4>
              <p className="text-sm text-slate-500 font-medium">Pesan layanan via dashboard atau aplikasi.</p>
            </div>

            <div className="bg-white p-6 rounded-3xl w-64 border border-[#e4ece3] shadow-sm relative z-10">
              <div className="w-16 h-16 bg-[#64965a] text-white rounded-full flex items-center justify-center text-2xl font-bold mx-auto mb-4 border-4 border-white shadow-md">2</div>
              <h4 className="text-lg font-bold text-slate-800 mb-2">Kurir Menjemput</h4>
              <p className="text-sm text-slate-500 font-medium">Fleet agent akan menjemput paket ke lokasi Anda.</p>
            </div>

            <div className="bg-white p-6 rounded-3xl w-64 border border-[#e4ece3] shadow-sm relative z-10">
              <div className="w-16 h-16 bg-[#64965a] text-white rounded-full flex items-center justify-center text-2xl font-bold mx-auto mb-4 border-4 border-white shadow-md">3</div>
              <h4 className="text-lg font-bold text-slate-800 mb-2">Transit & Sortir</h4>
              <p className="text-sm text-slate-500 font-medium">Proses sortir super cepat di warehouse kami.</p>
            </div>

            <div className="bg-white p-6 rounded-3xl w-64 border border-[#e4ece3] shadow-sm relative z-10">
              <div className="w-16 h-16 bg-yellow-500 text-white rounded-full flex items-center justify-center text-2xl font-bold mx-auto mb-4 border-4 border-white shadow-md">4</div>
              <h4 className="text-lg font-bold text-slate-800 mb-2">Paket Diterima</h4>
              <p className="text-sm text-slate-500 font-medium">Paket sampai tujuan lengkap dengan ePOD.</p>
            </div>
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-20 px-4">
        <div className="max-w-5xl mx-auto bg-slate-900 rounded-[3rem] p-12 md:p-20 text-center relative overflow-hidden shadow-2xl">
          <div className="absolute top-0 right-0 w-64 h-64 bg-[#64965a]/20 rounded-full blur-3xl"></div>
          <div className="absolute bottom-0 left-0 w-64 h-64 bg-indigo-500/20 rounded-full blur-3xl"></div>
          
          <h2 className="text-3xl md:text-5xl font-extrabold text-white mb-6 relative z-10">Siap Untuk Mengirim Paket?</h2>
          <p className="text-lg text-slate-300 mb-10 max-w-2xl mx-auto font-medium relative z-10">
            Bergabunglah dengan ribuan pebisnis dan individu yang telah mempercayakan pengirimannya kepada Buroqet.
          </p>
          <div className="flex flex-col sm:flex-row justify-center gap-4 relative z-10">
            <Link to="/customer/register" className="px-8 py-4 bg-[#64965a] hover:bg-[#53804a] text-white font-extrabold rounded-2xl shadow-xl hover:-translate-y-1 transition-all">
              Daftar Sekarang — Gratis
            </Link>
            <Link to="/pricing" className="px-8 py-4 bg-white/10 hover:bg-white/20 text-white font-extrabold rounded-2xl transition-all">
              Cek Ongkos Kirim
            </Link>
          </div>
        </div>
      </section>
    </div>
  );
}
