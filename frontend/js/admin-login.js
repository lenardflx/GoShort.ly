
const storedHash = "$2a$12$HAa3r.gW0C3R5L68aX9LBu.0JMDBjWugbOHqnFrizITGou/RVAefG"; // for "admin123"

document.getElementById("admin-login-form").addEventListener("submit", async function (e) {
        e.preventDefault();
        const password = document.getElementById("admin-password").value;
        const errorText = document.getElementById("login-error");
        const match = await bcrypt.compare(password, storedHash);

        if (match) {
            document.getElementById("auth-section").remove();
            fetch("admin-panel.html")
                .then((res) => res.text())
                .then((html) => {
                    document.getElementById("admin-placeholder").innerHTML = html;
                    attachSectionToggles();
                });
        } else {
            errorText.textContent = "Incorrect password.";
        }
    });