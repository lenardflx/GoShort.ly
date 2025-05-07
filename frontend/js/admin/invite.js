export function attachInviteHandler() {
    const btn = document.getElementById("send-invite-btn");
    const emailInput = document.getElementById("invite-email");
    const result = document.getElementById("invite-result");

    if (!btn || !emailInput || !result) return;

    btn.addEventListener("click", () => {
        const email = emailInput.value.trim();
        if (!email) return alert("Please enter an email.");

        btn.disabled = true;
        btn.textContent = "Sending...";

        fetch("http://localhost:8080/api/invite", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ email })
        })
            .then(res => {
                if (!res.ok) return res.text().then(err => { throw new Error(err); });
                return res.json();
            })
            .then(data => {
                result.innerHTML = `✅ Invite created. <br> <code>${data.invite}</code>`;
                emailInput.value = "";
            })
            .catch(err => {
                result.innerHTML = `❌ Failed: ${err.message}`;
            })
            .finally(() => {
                btn.disabled = false;
                btn.textContent = "Send Invite Link";
            });
    });
}
