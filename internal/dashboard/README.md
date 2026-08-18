# Dashboard Backend Modules

Setiap folder role berisi handler backend yang terisolasi. Router utama wajib memasang middleware autentikasi dan role sebelum mendaftarkan handler. Modul tidak boleh mempercayai role dari browser; role harus berasal dari session/token yang diverifikasi server.
