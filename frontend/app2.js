const table = document.getElementById("table");

let tbody;
table.childNodes.forEach((n) => {
    if (n.tagName == "TBODY") {
        tbody = n;
    }
});

fetch("http://localhost:8080/admin-ekran")
    .then((res) => res.json())
    .then((values) => {
        values.forEach((val) => {
            console.log();
            let row = `
                <tr>
                    <td>${val.id}</td>
                    <td>${val.marka}</td>
                    <td>${val.model}</td>
                    <td>${val.fiyat}</td>
                    <td>${val.stock}</td>
                    <td>${val.taleb_count.Int32}</td>
                </tr>
            `;
            tbody.insertAdjacentHTML("beforeend", row);
        });
    })
    .catch((err) => {
        console.log(err);
    });

tbody.addEventListener("click", function (e) {
    if (e.target.tagName === "BUTTON") {
        let cihazId = e.target.dataset.cihazId;
        if (cihazId != null) {
            fetch("http://localhost:8080/indirim-taleb", {
                method: "POST",
                mode: "no-cors",
                body: JSON.stringify({ user: "kudim", cihaz_id: +cihazId }),
            })
                .then((res) => {})
                .catch((err) => {
                    console.log(err);
                });
        }
    }
});
