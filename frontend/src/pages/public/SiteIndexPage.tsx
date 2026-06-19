import { Link } from 'react-router-dom';

export default function SiteIndexPage() {
  const sections = [
    {
      title: "Public Domain",
      links: [
        { path: "/", label: "Landing Page" },
        { path: "/tracking", label: "Cek Resi (Tracking)" },
        { path: "/pricing", label: "Cek Ongkir (Pricing)" },
      ]
    },
    {
      title: "Authentication Domain",
      links: [
        { path: "/customer/login", label: "Customer Login" },
        { path: "/customer/register", label: "Customer Register" },
        { path: "/courier/login", label: "Courier Login" },
        { path: "/admin/login", label: "Admin Login" },
      ]
    },
    {
      title: "Customer Domain (Protected)",
      links: [
        { path: "/customer/dashboard", label: "Customer Dashboard" },
        { path: "/customer/orders", label: "Manajemen Order" },
        { path: "/customer/orders/create", label: "Buat Order Baru" },
        { path: "/customer/profile", label: "Profil Saya" },
      ]
    },
    {
      title: "Courier Domain (Protected)",
      links: [
        { path: "/courier/dashboard", label: "Courier Dashboard" },
        { path: "/courier/assignment", label: "Pickup / Delivery" },
        { path: "/courier/epod", label: "Upload ePOD" },
        { path: "/courier/history", label: "Riwayat Pengiriman" },
      ]
    },
    {
      title: "Admin Domain (Protected)",
      links: [
        { path: "/admin/dashboard", label: "Admin Dashboard" },
        { path: "/admin/orders", label: "Manajemen Order (Tabel)" },
        { path: "/admin/warehouse", label: "Warehouse Management" },
        { path: "/admin/dispatch", label: "Dispatch Fleet" },
        { path: "/admin/pricing", label: "Pricing & Routing" },
        { path: "/admin/settlement", label: "Settlement (COD)" },
        { path: "/admin/tracking", label: "Tracking Monitor" },
      ]
    }
  ];

  return (
    <div className="min-h-screen bg-slate-50 py-12 px-4 sm:px-6 lg:px-8 font-sans">
      <div className="max-w-4xl mx-auto">
        <div className="text-center mb-12">
          <h1 className="text-4xl font-extrabold text-slate-900 mb-4">🗺️ Sitemap / Dev Index</h1>
          <p className="text-lg text-slate-600">
            Daftar seluruh halaman di aplikasi Buroqet untuk mempermudah navigasi selama proses development.
          </p>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {sections.map((section, idx) => (
            <div key={idx} className="bg-white rounded-2xl shadow-sm border border-slate-200 p-6 hover:shadow-md transition-shadow">
              <h2 className="text-xl font-bold text-slate-800 border-b-2 border-slate-100 pb-3 mb-4">
                {section.title}
              </h2>
              <ul className="space-y-3">
                {section.links.map((link, lidx) => (
                  <li key={lidx}>
                    <Link 
                      to={link.path}
                      className="group flex items-center text-slate-600 hover:text-[#64965a] font-medium transition-colors"
                    >
                      <span className="w-6 h-6 rounded-full bg-slate-100 group-hover:bg-[#64965a]/10 flex items-center justify-center text-xs mr-3 transition-colors">
                        ↗
                      </span>
                      {link.label}
                    </Link>
                  </li>
                ))}
              </ul>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
