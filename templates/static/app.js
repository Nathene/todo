document.addEventListener("DOMContentLoaded", () => {
    // Handle checkbox change events
    document.querySelectorAll('input[type="checkbox"]').forEach((checkbox) => {
        checkbox.addEventListener("change", async (event) => {
            const subtaskId = event.target.dataset.subtaskId;
            const done = event.target.checked;

            try {
                const response = await fetch(`/subtasks/${subtaskId}/toggle`, {
                    method: "POST",
                    body: JSON.stringify({ done }),
                    headers: {
                        "Content-Type": "application/json",
                    },
                });

                if (!response.ok) {
                    throw new Error("Failed to update subtask status");
                }

                // Optionally update UI styles
                const parentLi = event.target.closest("li");
                if (done) {
                    parentLi.classList.add("bg-success", "text-decoration-line-through");
                } else {
                    parentLi.classList.remove("bg-success", "text-decoration-line-through");
                }
            } catch (error) {
                console.error("Error updating subtask:", error);
                alert("Failed to update subtask status");
            }
        });
    });
});
