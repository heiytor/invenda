import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faEnvelope } from "@fortawesome/free-solid-svg-icons";
import { FieldError, UseFormRegisterReturn } from "react-hook-form";

type InputEmailProps = {
    email: UseFormRegisterReturn; // Atualize a prop para aceitar UseFormRegisterReturn
    errors?: FieldError | undefined;
};

export function InputEmail({ email, errors }: InputEmailProps) {
    return (
        <div className="relative">
            <input
                id="email"
                type="text"
                {...email} // Use o objeto de registro retornado por register
                placeholder="Email"
                className="border rounded p-2 w-full pl-8"
            />
            <div className="absolute inset-y-0 left-0 flex items-center pl-2">
                <FontAwesomeIcon icon={faEnvelope} className="text-gray-400" />
            </div>
            {errors && (
                <span className="text-red-500 text-sm mt-1">{errors.message}</span>
            )}
        </div>
    );
}
