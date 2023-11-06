/// <reference path="../libs/tabulator-tables/tabulator.min.js"/>
let fake = [
    {
        name: "foo",
        gender: "male",
        col: "black",
    },
];

var table = new Tabulator("#tbRank", {
    data: fake,
    height: "311px",
    columns: [
        { title: "Name", field: "name" },
        { title: "Progress", field: "progress", sorter: "number" },
        { title: "Gender", field: "gender" },
        { title: "Rating", field: "rating" },
        { title: "Favourite Color", field: "col" },
        { title: "Date Of Birth", field: "dob", hozAlign: "center" },
    ],
});