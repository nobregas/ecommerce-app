import { Colors } from "@/mobile/constants/Colors";
import React from "react";
import { StyleSheet, TextInput, TextInputProps, Text } from "react-native";

interface InputFieldProps extends TextInputProps {
}
const inputField = ({ ...props }: InputFieldProps) => {
    return (
        <>
            <TextInput
                {...props}
                style={styles.inputField}
            />
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

})
