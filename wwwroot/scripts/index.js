/// <reference path="../libs/tabulator-tables/tabulator.min.js"/>
let top10Cuisines = [];
window.addEventListener("DOMContentLoaded", function(){
    /** @type {HTMLTextAreaElement} */
    let top10Viewbag = document.querySelector(".viewbag#top10Cuisines");
    top10Cuisines = JSON.parse(top10Viewbag.value);

    let tbTop10Cuisines = new Tabulator("#tbTop10Cuisines", {
        height: "240px",
        data: top10Cuisines,
        columns: [
            { title: "餐點名稱", field: "cuisine_name" },
            { title: "單價", field: "unit_price"},
            { title: "套餐", field: "is_one_set"},
            { title: "日期", field: "last_order_date" },
            { title: "評分", field: "review" },
            { title: "餐廳", field: "restaurant" },
            { title: "地點", field: "address" },
            { title: "備註", field: "remark" },
        ],
    });

});
