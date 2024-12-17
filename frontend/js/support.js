const purch_btn = document.getElementById("support");
const modal3 = document.getElementById("modal3");
const closeModal3 = document.getElementById("closeModal3");

const purch_card_btn = document.getElementById("purch-card-btn");
const modal4 = document.getElementById("modal4");
const closeModal4 = document.getElementById("closeModal4");

purch_btn.onclick = ()=>{
    modal3.style.display = "block";
}

closeModal3.onclick = ()=>{
    modal3.style.display = "none";
}

purch_card_btn.onclick = ()=>{
    modal4.style.display = "block";
    modal3.style.display = "none";
}

closeModal4.onclick = ()=>{
    modal4.style.display = "none";
}

window.onclick = function(event) {
    if (event.target === modal3) {
        modal3.style.display = "none";
    }
    if (event.target === modal4) {
        modal4.style.display = "none";
    }
}