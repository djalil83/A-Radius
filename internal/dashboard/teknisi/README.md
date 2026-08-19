# Teknisi Dashboard Backend

Modul backend untuk role `teknisi`. Handler di sini harus didaftarkan melalui router aplikasi dan dilindungi middleware RBAC.

Endpoint yang spesifik untuk role `teknisi` diletakkan di modul ini. Endpoint lintas-role seperti Subscription Profile tetap menggunakan modul domain terkait di `internal/subscriptionprofile`.
