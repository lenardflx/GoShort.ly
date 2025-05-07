document.getElementById("send-test-email").addEventListener("click", () => {
    const form = document.getElementById("admin-settings-form");
    const data = new FormData(form);
    const json = Object.fromEntries(data.entries());

    json.smtp_port = parseInt(json.smtp_port || "587");

    console.log("Sending test email with data:", json);

    fetch("http://localhost:8080/api/send-test-email", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(json)
    })
        .then(res => {
            if (!res.ok) return res.text().then(err => {
                throw new Error(err);
            });
            return res.json();
        })
        .then(() => {
            alert("Test email sent successfully!");
        })
        .catch(err => {
            alert("Failed to send: " + err.message);
        });
});
