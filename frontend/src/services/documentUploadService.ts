import { getToken } from "@clerk/react";

import { API_URL } from "@/api/client";

export async function uploadDocument(
    file: File,
) {

    const formData = new FormData();

    formData.append(
        "file",
        file,
    );

    const token = await getToken();

    const response = await fetch(
        `${API_URL}/documents/upload`,
        {
            method: "POST",
            body: formData,
            headers: token ? { Authorization: `Bearer ${token}` } : undefined,
        },
    );

    if (!response.ok) {

        const errorMessage = await response.text();

        throw new Error(
            errorMessage,
        );

        // throw new Error(
        //     `Upload failed (${response.status})`,
        // );
    }

    return response.json();
}
