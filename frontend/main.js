const btn_login = document.getElementById("btn-login");
const modal = document.getElementById("modal");
const reg_btn = document.getElementById("reg-btn");
const closeModal = document.getElementById("closeModal");
const closeModal2 = document.getElementById("closeModal2");

btn_login.onclick = ()=>{
    modal.style.display = "block";
}

reg_btn.onclick = ()=>{
    modal.style.display = "none";
    modal2.style.display = "block";
}

closeModal.onclick = ()=>{
    modal.style.display = "none";
}

closeModal2.onclick = ()=>{
    modal2.style.display = "none";
}

window.onclick = function(event) {
    if (event.target === modal) {
        modal.style.display = "none";
    }
    if (event.target === modal2) {
        modal2.style.display = "none";
    }
}