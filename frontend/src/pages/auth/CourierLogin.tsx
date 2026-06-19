import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '../../store/auth.store';

export default function CourierLogin() {
  const { login, isLoading, error } = useAuthStore();
  const navigate = useNavigate();
  
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await login(email, password);
      navigate('/courier/dashboard');
    } catch (err) {
      // Error handled by store
    }
  };

  return (
    <div className="bg-white shadow-2xl rounded-3xl overflow-hidden border border-[#e4ece3]">
      <div className="p-8 md:p-10">
        <div className="text-center mb-8">
          <span className="text-5xl text-[#64965a] mb-4 block">🛵</span>
          <h2 className="text-3xl font-extrabold text-slate-800">Fleet Agent Login</h2>
          <p className="text-slate-500 text-sm mt-2">Buroqet Courier Operations</p>
        </div>
        
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label className="block text-sm font-bold text-slate-700 mb-2 uppercase tracking-wide" htmlFor="email">Email Address</label>
            <input
              id="email"
              type="email"
              className="w-full px-5 py-4 rounded-2xl border border-[#e4ece3] focus:outline-none focus:ring-4 focus:ring-[#64965a]/20 focus:border-[#64965a] transition-all bg-slate-50 text-slate-900 font-medium"
              placeholder="courier@buroqet.com"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              required
            />
          </div>

          <div>
            <div className="flex justify-between items-center mb-2">
              <label className="block text-sm font-bold text-slate-700 uppercase tracking-wide" htmlFor="password">Password</label>
              <span className="text-sm font-semibold text-[#64965a] hover:underline cursor-pointer">Lupa Password?</span>
            </div>
            <input
              id="password"
              type="password"
              className="w-full px-5 py-4 rounded-2xl border border-[#e4ece3] focus:outline-none focus:ring-4 focus:ring-[#64965a]/20 focus:border-[#64965a] transition-all bg-slate-50 text-slate-900 font-medium"
              placeholder="Enter your password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
          </div>

          {error && (
            <div className="bg-red-50 text-red-600 p-4 rounded-2xl text-sm font-bold border border-red-100 flex items-start gap-3">
              <span className="text-xl">⚠️</span> {error}
            </div>
          )}

          <button
            type="submit"
            className="w-full py-4 px-6 bg-[#64965a] hover:bg-[#53804a] text-white font-extrabold text-lg rounded-2xl shadow-lg hover:shadow-xl transform hover:-translate-y-1 transition-all duration-300 focus:outline-none focus:ring-4 focus:ring-[#64965a]/50 disabled:opacity-70 disabled:cursor-not-allowed"
            disabled={isLoading}
          >
            {isLoading ? 'Authenticating...' : 'Sign In'}
          </button>
        </form>

        <div className="mt-8 text-center">
          <p className="text-slate-600 font-medium">
            Belum terdaftar?{' '}
            <Link to="/courier/register" className="text-[#64965a] font-bold hover:underline">
              Sign Up
            </Link>
          </p>
        </div>
      </div>
    </div>
  );
}
