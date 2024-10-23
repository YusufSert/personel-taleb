const table = document.getElementById("table");

let tbody;
table.childNodes.forEach((n) => {
    if (n.tagName == "TBODY") {
        tbody = n;
    }
});

fetch("http://localhost:8080/indirim-taleb")
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
                    <td>${val.taleb_count}/${
                val.required_taleb_count.Int32
            }</td>
            <td><button id="btn-taleb" data-cihaz-id=${val.id}${
                val.taleb_state.String == "active" ? " disabled" : ""
            } >${
                val.taleb_count == 0 ? "İndirim İste" : val.taleb_count
            }</button>
            </td>
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
/*
// select pc.id, count(pt.id), (select gerekli_taleb_count from personeal_indirim_taleb_realization where
 taleb_state IN ('pending', 'active') and cihaz_id = pc.id) from personel_cihazlar pc left join
  (select id,  cihaz_id from personel_taleb where evaulated = 0 )pt on(pc.id = pt.cihaz_id) group by pc.id;
*/

// select pc.id, count(pt.id), ptr.cihaz_id from personel_cihazlar pc left join (select id,  cihaz_id from personel_taleb where evaulated = 0 )pt on(pc.id = pt.cihaz_id) left join (select * from personeal_indirim_taleb_realization where taleb_state IN ('pending', 'active')) ptr on(ptr.cihaz_id = pc.id) group by pc.id, ptr.cihaz_id;

// select pc.id, pc.marka, pc.model, pc.fiyat, pc.stock, count(pt.id), ptr.gerekli_taleb_count, ptr.taleb_state from personel_cihazlar pc left join (select id,  cihaz_id from personel_taleb where evaulated = 0)pt on(pc.id = pt.cihaz_id) left join (select * from personeal_indirim_taleb_realization where taleb_state IN ('pending', 'active')) ptr on(ptr.cihaz_id = pc.id) group by pc.id, ptr.cihaz_id, ptr.taleb_state, ptr.gerekli_taleb_count;

// select pc.id, pc.marka, pc.model, pc.fiyat, pc.stock, count(pt.id) as taleb_count, ptr.gerekli_taleb_count, ptr.taleb_state from personel_cihazlar pc left join (select id,  cihaz_id from personel_taleb where evaulated = 0)pt on(pc.id = pt.cihaz_id) left join (select * from personeal_indirim_taleb_realization where taleb_state IN ('pending', 'active')) ptr on(ptr.cihaz_id = pc.id) group by pc.id, ptr.cihaz_id, ptr.taleb_state, ptr.gerekli_taleb_count;
