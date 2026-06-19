import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuthStore } from '../../store/auth.store';
import './LoginPage.css';

export default function RegisterPage() {
  const { register, isLoading, error } = useAuthStore();
  const navigate = useNavigate();
  
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [role, setRole] = useState('pelanggan');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    try {
      await register(name, email, password, role);
      alert('Registrasi berhasil! Silakan login.');
      navigate('/login');
    } catch (err) {
      // error sudah dihandle oleh store
    }
  };

  return (
    <form className="login-form" onSubmit={handleSubmit} style={{ marginTop: '20px' }}>
      <h2 style={{ textAlign: 'center', marginBottom: '20px', fontSize: '24px', fontWeight: 'bold' }}>Buat Akun</h2>
      
      <div className="form-group">
        <label className="form-label" htmlFor="name">Nama Lengkap</label>
        <input
          id="name"
          type="text"
          className="input"
          placeholder="Nama Anda"
          value={name}
          onChange={(e) => setName(e.target.value)}
          required
        />
      </div>

      <div className="form-group">
        <label className="form-label" htmlFor="email">Email</label>
        <input
          id="email"
          type="email"
          className="input"
          placeholder="you@buroqet.com"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          required
        />
      </div>

      <div className="form-group">
        <label className="form-label" htmlFor="password">Password</label>
        <input
          id="password"
          type="password"
          className="input"
          placeholder="Min. 6 karakter"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
          minLength={6}
        />
      </div>

      <div className="form-group">
        <label className="form-label" htmlFor="role">Daftar Sebagai</label>
        <select 
          id="role" 
          className="input" 
          value={role} 
          onChange={(e) => setRole(e.target.value)}
          required
        >
          <option value="pelanggan">Pelanggan</option>
          <option value="kurir">Kurir</option>
          <option value="admin">Admin</option>
        </select>
      </div>

      {error && <div className="error-msg">⚠️ {error}</div>}

      <button
        type="submit"
        className="btn btn-primary login-btn"
        disabled={isLoading}
        style={{ marginTop: '10px' }}
      >
        {isLoading ? 'Mendaftar...' : 'Daftar'}
      </button>

      <div style={{ textAlign: 'center', marginTop: '16px', fontSize: '13px', color: 'var(--text-secondary)' }}>
        Sudah punya akun? <Link to="/login" style={{ color: 'var(--primary)', fontWeight: 600 }}>Login di sini</Link>
      </div>
    </form>
  );
}
