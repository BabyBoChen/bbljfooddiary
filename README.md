# BBLJ食記

「BBLJ食記」記錄了BBLJ團隊平日外食的地點、餐點價格以及評價，供有興趣的朋友們參考。

「BBLJ食記」是以Golang語言、基於fiber框架編寫而成的一個網站作品，經包裝成Docker映像後部署於微軟的Azure平台上。

本站採用了Neon的Postgres資料庫以及Dropbox雲端硬碟作為後端的資料儲存解決方案，實作方法請參考本專案的dropboxservice.go以及pgdbcontext套件( https://github.com/BabyBoChen/pgdbcontext )。

網站連結：
https://bbljfooddiary.azurewebsites.net/
