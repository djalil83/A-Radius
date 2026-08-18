package shared

// DashboardContract mendefinisikan identitas dan hak akses dasar setiap dashboard.
type DashboardContract struct {
	Role        string
	DisplayName string
	BasePath    string
	Permissions []string
}

// Contracts adalah sumber metadata dashboard untuk dokumentasi dan validasi route.
var Contracts = map[string]DashboardContract{
	"developer":     {Role: "developer", DisplayName: "Developer Security Dashboard", BasePath: "/dashboard/developer", Permissions: []string{"system:read", "system:write", "security:read", "security:scan", "threat:read", "credential:read", "credential:rotate", "dependency:read", "dependency:fix", "api:read", "database:read", "audit:read", "code:read", "code:write", "preview:read", "approval:read", "approval:decide", "deployment:read", "deployment:run", "deployment:rollback"}},
	"administrator": {Role: "administrator", DisplayName: "Administrator Dashboard", BasePath: "/dashboard/administrator", Permissions: []string{"profile:read", "profile:write", "user:manage", "audit:read"}},
	"teknisi":       {Role: "teknisi", DisplayName: "Teknisi Dashboard", BasePath: "/dashboard/teknisi", Permissions: []string{"ticket:read", "ticket:update", "network:read"}},
	"reseller":      {Role: "reseller", DisplayName: "Reseller Dashboard", BasePath: "/dashboard/reseller", Permissions: []string{"customer:read", "subscription:read", "subscription:create"}},
	"pelanggan":     {Role: "pelanggan", DisplayName: "Pelanggan Dashboard", BasePath: "/dashboard/pelanggan", Permissions: []string{"profile:read", "subscription:read", "ticket:create"}},
}
