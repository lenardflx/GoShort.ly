export function attachSmtpTestHandler() {
    const btn = document.getElementById("send-test-email");
    if (!btn) return console.warn("SMTP button not found");

    btn.addEventListener("click", () => {
        const form = document.getElementById("admin-settings-form");
        const data = new FormData(form);
        const json = Object.fromEntries(data.entries());

        console.table(json);

        // Convert smtp_port to number
        json.smtp_port = parseInt(json.smtp_port || "587");

        // Extract test recipient address
        const recipient = json.smtp_test_to;
        if (!recipient) return alert("Please enter a test recipient email address.");

        delete json.smtp_test_to; // remove from config payload

        const payload = {
            ...json,
            smtp_test_to: recipient,
        };

        // UX: disable button during request
        btn.disabled = true;
        btn.textContent = "Sending...";

        fetch("http://localhost:8080/api/send-test-email", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify(payload)
        })
            .then(res => {
                if (!res.ok) return res.text().then(err => { throw new Error(err); });
                return res.json();
            })
            .then(() => alert("✅ Test email sent!"))
            .catch(err => alert("❌ Failed: " + err.message))
            .finally(() => {
                btn.disabled = false;
                btn.textContent = "Send Test Email";
            });
    });
}
