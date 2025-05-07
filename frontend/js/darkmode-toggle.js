const toggle = document.getElementById("darkmode-toggle");
const dark = localStorage.getItem("theme") === "dark";
if (dark) document.body.classList.add("dark");

toggle.addEventListener("click", () => {
    document.body.classList.toggle("dark");
    localStorage.setItem("theme", document.body.classList.contains("dark") ? "dark" : "light");
});