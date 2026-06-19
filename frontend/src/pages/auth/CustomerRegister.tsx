import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '../../store/auth.store';

import apiClient from '../../services/api-client';
import { API_ENDPOINTS } from '../../utils/api-config';

export default function RegisterPage() {
  const navigate = useNavigate();
  const [localLoading, setLocalLoading] = useState(false);
  const [localError, setLocalError] = useState('');
  
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setLocalLoading(true);
    setLocalError('');
    try {
      // Secara otomatis mendaftar sebagai pelanggan
      await apiClient.post(API_ENDPOINTS.auth.register, { name, email, password, role: 'pelanggan' });
      alert('Registrasi berhasil! Silakan login.');
      navigate('/customer/login');
    } catch (err: any) {
      setLocalError(err.response?.data?.error || 'Gagal registrasi');
    } finally {
      setLocalLoading(false);
    }
  };

  return (
    <div className="bg-white shadow-2xl rounded-3xl overflow-hidden border border-[#e4ece3]">
      {/* Tab Nav */}
      <div className="flex border-b border-[#e4ece3]">
        <Link to="/customer/login" className="flex-1 text-center py-5 text-slate-500 hover:text-slate-700 hover:bg-slate-50 transition-colors font-bold">
          Sign In
        </Link>
        <div className="flex-1 text-center py-5 text-[#64965a] border-b-4 border-[#64965a] font-bold bg-[#64965a]/5">
          Sign Up
        </div>
      </div>

      <div className="p-8 md:p-10">
        <div className="text-center mb-8">
          <span className="text-5xl text-[#64965a] mb-4 block">✨</span>
          <h2 className="text-3xl font-extrabold text-slate-800">Buat Akun Baru</h2>
          <p className="text-slate-500 text-sm mt-2">Daftar untuk mulai mengirim paket</p>
        </div>
        
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-bold text-slate-700 mb-2 uppercase tracking-wide" htmlFor="name">Nama Lengkap</label>
            <input
              id="name"
              type="text"
              className="w-full px-5 py-4 rounded-2xl border border-[#e4ece3] focus:outline-none focus:ring-4 focus:ring-[#64965a]/20 focus:border-[#64965a] transition-all bg-slate-50 text-slate-900 font-medium"
              placeholder="Nama Anda"
              value={name}
              onChange={(e) => setName(e.target.value)}
              required
            />
          </div>

          <div>
            <label className="block text-sm font-bold text-slate-700 mb-2 uppercase tracking-wide" htmlFor="email">Email Address</label>
            <input
              id="email"
              type="email"
              className="w-full px-5 py-4 rounded-2xl border border-[#e4ece3] focus:outline-none focus:ring-4 focus:ring-[#64965a]/20 focus:border-[#64965a] transition-all bg-slate-50 text-slate-900 font-medium"
              placeholder="you@buroqet.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          <div>
            <label className="block text-sm font-bold text-slate-700 mb-2 uppercase tracking-wide" htmlFor="password">Password</label>
            <input
              id="password"
              type="password"
              className="w-full px-5 py-4 rounded-2xl border border-[#e4ece3] focus:outline-none focus:ring-4 focus:ring-[#64965a]/20 focus:border-[#64965a] transition-all bg-slate-50 text-slate-900 font-medium"
              placeholder="Min. 6 karakter"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
              minLength={6}
            />
          </div>

          {localError && (
            <div className="bg-red-50 text-red-600 p-4 rounded-2xl text-sm font-bold border border-red-100 flex items-start gap-3">
              <span className="text-xl">⚠️</span> {localError}
            </div>
          )}

          <button
            type="submit"
            className="w-full py-4 px-6 bg-[#64965a] hover:bg-[#53804a] text-white font-extrabold text-lg rounded-2xl shadow-lg hover:shadow-xl transform hover:-translate-y-1 transition-all duration-300 focus:outline-none focus:ring-4 focus:ring-[#64965a]/50 disabled:opacity-70 disabled:cursor-not-allowed"
            disabled={localLoading}
          >
            {localLoading ? 'Memproses...' : 'Daftar Sekarang'}
          </button>
        </form>
      </div>
    </div>
  );
}
