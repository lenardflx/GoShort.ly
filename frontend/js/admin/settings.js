export function attachAdminFormSave() {
    const form = document.getElementById("admin-settings-form");
    const feedback = document.getElementById("form-feedback");

    loadAdminConfig();

    form.addEventListener("submit", (e) => {
        e.preventDefault();
        const data = new FormData(form);
        const json = Object.fromEntries(data.entries());

        json.max_url_length = +json.max_url_length || 0;
        json.default_user_limit = +json.default_user_limit || 0;
        json.default_expiration = +json.default_expiration || 0;
        json.failed_login_limit = +json.failed_login_limit || 0;
        json.user_override_expiration = !!data.get("user_override_expiration");
        json.allow_anonymous = !!data.get("allow_anonymous");

        fetch("http://localhost:8080/api/admin-config/save", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(json),
        })
            .then(res => {
                if (!res.ok) throw new Error("Save failed");
                return res.json();
            })
            .then(() => {
                feedback.textContent = "✅ Settings saved.";
                feedback.style.color = "green";
            })
            .catch(err => {
                feedback.textContent = "❌ Save failed: " + err.message;
                feedback.style.color = "red";
            });
    });
}

function loadAdminConfig() {
    fetch("http://localhost:8080/api/admin-config")
        .then(res => res.json())
        .then(data => {
            for (const [key, value] of Object.entries(data)) {
                const el = document.querySelector(`[name="${key}"]`);
                if (!el) continue;

                if (el.type === "checkbox") el.checked = Boolean(value);
                else el.value = value;
            }
        });
}
