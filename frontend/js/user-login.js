const toggle = document.getElementById("darkmode-toggle");
const dark = localStorage.getItem("theme") === "dark";
if (dark) document.body.classList.add("dark");

toggle.addEventListener("click", () => {
    document.body.classList.toggle("dark");
    localStorage.setItem("theme", document.body.classList.contains("dark") ? "dark" : "light");
});

const loginSection = document.getElementById("login-section");
const userPlaceholder = document.getElementById("user-placeholder");
const userLoginForm = document.getElementById("user-login-form");

function setupDashboard(username) {
    document.getElementById("display-name").textContent = username;
    document.getElementById("settings-username").textContent = username;

    // Logout
    document.getElementById("logout-btn").addEventListener("click", () => {
        localStorage.removeItem("username");
        location.reload();
    });

    // Search filter
    const searchInput = document.getElementById("search-links");
    searchInput.addEventListener("input", () => {
        const filter = searchInput.value.toLowerCase();
        document.querySelectorAll("#link-list li").forEach(li => {
            li.style.display = li.textContent.toLowerCase().includes(filter) ? "" : "none";
        });
    });

    // QR Code Generator (placeholder)
    document.getElementById("link-form").addEventListener("submit", (e) => {
        e.preventDefault();
        const url = document.getElementById("target-url").value.trim();
        const slug = document.getElementById("custom-slug").value.trim();
        if (!url) return;

        const result = `✅ Created: <strong>${slug || "(auto)"}</strong> → ${url}<br/>[QR Code Placeholder]`;
        document.getElementById("qr-result").innerHTML = result;
    });
}

function showDashboard(username) {
    fetch("user-dashboard.html")
        .then(res => res.text())
        .then(html => {
            loginSection.remove();
            userPlaceholder.innerHTML = html;
            setupDashboard(username);
        });
}

// Automatically show dashboard if user is cached
const cachedUser = localStorage.getItem("username");
if (cachedUser) {
    showDashboard(cachedUser);
}

// Login form handler
userLoginForm.addEventListener("submit", (e) => {
    e.preventDefault();
    const username = document.getElementById("username").value.trim();
    if (!username) return;
    localStorage.setItem("username", username);
    showDashboard(username); // ✅ Call it here
});
