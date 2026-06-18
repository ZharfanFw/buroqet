

interface StatCardProps {
  title: string;
  value: string | number;
  subtitle?: string;
  trend?: 'up' | 'down' | 'neutral';
  color?: 'primary' | 'success' | 'warning' | 'danger' | 'info' | 'default';
  className?: string;
}

export function StatCard({ title, value, subtitle, trend, color = 'default', className = '' }: StatCardProps) {
  // Map our custom CSS variable names or Tailwind equivalents
  const valueColors = {
    primary: 'text-[color:var(--primary)]',
    success: 'text-[color:var(--success)]',
    warning: 'text-[color:var(--warning)]',
    danger: 'text-[color:var(--danger)]',
    info: 'text-[color:var(--info)]',
    default: 'text-[color:var(--text-primary)]'
  };

  const subtitleColors = {
    primary: 'text-[color:var(--primary-light)]',
    success: 'text-[color:var(--success)]',
    warning: 'text-[color:var(--text-secondary)]',
    danger: 'text-[color:var(--danger)]',
    info: 'text-[color:var(--text-secondary)]',
    default: 'text-[color:var(--success)]' // Default trend is green
  };

  const trendIcon = trend === 'up' ? '↑' : trend === 'down' ? '↓' : '';

  return (
    <div className={`card flex flex-col justify-center p-5 ${className}`}>
      <div className="text-[color:var(--text-muted)] text-[13px] font-semibold">{title}</div>
      <div className={`text-[32px] font-bold mt-2 ${valueColors[color]}`}>{value}</div>
      {subtitle && (
        <div className={`text-[12px] mt-1 ${subtitleColors[color]}`}>
          {trendIcon} {subtitle}
        </div>
      )}
    </div>
  );
}
