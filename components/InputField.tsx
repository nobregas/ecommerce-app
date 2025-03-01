import { Colors } from "@/constants/Colors";
import React from "react";
import { StyleSheet, TextInput, TextInputProps } from "react-native";

type Props = {} & TextInputProps

const inputField = (props: Props) => {
    return (
        <TextInput
            {...props}
            style={styles.inputField}
        />
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
    }
})
