import { Colors } from "@/constants/Colors";
import React from "react";
import { StyleSheet, TextInput, TextInputProps, Text } from "react-native";

interface InputFieldProps extends TextInputProps {
    error?: string
}
const inputField = ({ error, ...props }: InputFieldProps) => {
    return (
        <>
            <TextInput
                {...props}
                style={[styles.inputField, error && styles.errorInput]}
            />
            {error && <Text style={styles.errorText}>{error}</Text>}
        </>
    )
}

export default inputField;

const styles = StyleSheet.create({
    inputField: {
        backgroundColor: Colors.white,
        paddingVertical: 12,
        paddingHorizontal: 18,
        alignSelf: 'stretch',
        borderRadius: 5,
        fontSize: 16,
        color: Colors.black,
        marginBottom: 20,
    },
    errorInput: {
        borderColor: Colors.red,
        borderWidth: 1,
    },
    errorText: {
        color: 'red',
        fontSize: 12,
        marginTop: -10,
        marginBottom: 15,
        alignSelf: 'flex-start',
    },
})
