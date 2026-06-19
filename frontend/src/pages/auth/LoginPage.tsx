import { useState } from 'react';
import type { FormEvent } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../store/auth.store';
import './LoginPage.css';

export default function LoginPage() {
  const { login, isLoading, error } = useAuthStore();
  const navigate = useNavigate();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    await login(email, password);
    if (useAuthStore.getState().isAuthenticated) {
      navigate('/tracking');
    }
  };

  return (
    <form className="login-form" onSubmit={handleSubmit}>
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
          placeholder="••••••••"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          required
        />
      </div>

      {error && <div className="error-msg">⚠️ {error}</div>}

      <button
        id="login-submit"
        type="submit"
        className="btn btn-primary login-btn"
        disabled={isLoading}
      >
        {isLoading ? 'Masuk...' : 'Masuk'}
      </button>
    </form>
  );
}
