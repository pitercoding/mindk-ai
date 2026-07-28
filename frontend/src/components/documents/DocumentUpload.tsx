import { useState } from "react";
import { uploadDocument } from "../../services/documentUploadService";

interface DocumentUploadProps {
    onUploadSuccess: () => void;
}

export default function DocumentUpload({
    onUploadSuccess,
}: DocumentUploadProps) {
    const [file, setFile] = useState<File | null>(null);
    const [uploading, setUploading] = useState(false);
    const [message, setMessage] = useState("");

    function handleFileChange(
        event: React.ChangeEvent<HTMLInputElement>
    ) {
        const selectedFile = event.target.files?.[0];

        if (selectedFile) {
            setFile(selectedFile);
            setMessage("");
        }
    }

    async function handleUpload() {
        if (!file) {
            setMessage("Please select a file first.");
            return;
        }

        try {
            setUploading(true);
            setMessage("");

            await uploadDocument(file);

            setMessage("Document uploaded successfully.");

            setFile(null);

            onUploadSuccess();

        } catch (error) {
            console.error(error);

            setMessage(
                "Error uploading document."
            );

        } finally {
            setUploading(false);
        }
    }

    return (
        <div className="document-upload">

            <h2>
                Upload Document
            </h2>


            <div className="document-upload-actions">

                <label className="file-upload-button">

                    Choose file

                    <input
                        type="file"
                        onChange={handleFileChange}
                        disabled={uploading}
                        hidden
                    />

                </label>


                <button
                    onClick={handleUpload}
                    disabled={uploading}
                >
                    {
                        uploading
                            ? "Uploading..."
                            : "Upload document"
                    }
                </button>

            </div>

            {file && (
                <p>
                    Selected file: {file.name}
                </p>
            )}

            {message && (
                <p>
                    {message}
                </p>
            )}

        </div>
    );
}
