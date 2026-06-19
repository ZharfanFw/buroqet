import { Outlet } from 'react-router-dom';

export default function AuthLayout() {
  return (
    <div className="min-h-screen flex items-center justify-center bg-[#f7faf7] py-12 px-4 sm:px-6 lg:px-8 relative overflow-hidden">
      {/* Decorative Background Elements */}
      <div className="absolute top-[-10%] left-[-10%] w-[40%] h-[40%] rounded-full bg-[#64965a]/10 blur-3xl pointer-events-none" />
      <div className="absolute bottom-[-10%] right-[-10%] w-[40%] h-[40%] rounded-full bg-[#79ae6f]/10 blur-3xl pointer-events-none" />
      
      <div className="max-w-md w-full z-10">
        <div className="text-center mb-8">
          <div className="flex justify-center items-center gap-3 mb-4">
            <span className="text-4xl filter drop-shadow-md">🚀</span>
            <h1 className="text-4xl font-extrabold text-[#64965a] tracking-tight">
              Buroqet
            </h1>
          </div>
          <p className="mt-2 text-sm text-slate-500 font-bold uppercase tracking-widest">
            Sistem Logistik Masa Depan
          </p>
        </div>
        
        <Outlet />
      </div>
    </div>
  );
}
