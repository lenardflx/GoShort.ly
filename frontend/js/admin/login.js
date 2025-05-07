import { fetchTextHTML, attachSectionToggles } from "./utils.js";
import { attachAdminFormSave } from "./settings.js";
import { attachSmtpTestHandler } from "./test-smtp.js";
import { attachInviteHandler } from "./invite.js";

document.getElementById("admin-login-form").addEventListener("submit", async function (e) {
    e.preventDefault();
    const errorText = document.getElementById("login-error");
    const password = document.getElementById("admin-password").value;

    try {
        const res = await fetch("http://localhost:8080/api/admin-config/hash");
        const { hash } = await res.json();
        const match = await bcrypt.compare(password, hash);

        if (!match) {
            errorText.textContent = "❌ Incorrect password.";
            return;
        }

        // Load panel
        document.getElementById("auth-section").remove();
        const html = await fetchTextHTML("admin-panel.html");
        document.getElementById("admin-placeholder").innerHTML = html;

        attachSectionToggles();
        attachSmtpTestHandler();
        attachAdminFormSave();
        attachInviteHandler();

    } catch (err) {
        errorText.textContent = "❌ Login error: " + err.message;
    }
});
