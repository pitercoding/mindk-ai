import { useState } from "react";

import { useToast } from "@/context/toast-context";

import { uploadDocument } from "../../services/documentUploadService";

interface DocumentUploadProps {
    onUploadSuccess: () => void;
}

export default function DocumentUpload({
    onUploadSuccess,
}: DocumentUploadProps) {
    const [file, setFile] = useState<File | null>(null);
    const [uploading, setUploading] = useState(false);

    const toast = useToast();

    function handleFileChange(
        event: React.ChangeEvent<HTMLInputElement>
    ) {
        const selectedFile = event.target.files?.[0];

        if (selectedFile) {
            setFile(selectedFile);
        }
    }

    async function handleUpload() {
        if (!file) {
            toast.warning("Please select a file first.");
            return;
        }

        try {
            setUploading(true);

            await uploadDocument(file);

            toast.success("Document uploaded successfully.");

            setFile(null);

            onUploadSuccess();

        } catch (error) {
            console.error(error);

            toast.error("Error uploading document.");

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

                <label className="btn btn-secondary">

                    Choose file

                    <input
                        type="file"
                        onChange={handleFileChange}
                        disabled={uploading}
                        hidden
                    />

                </label>


                <button
                    className="btn btn-primary"
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

        </div>
    );
}
