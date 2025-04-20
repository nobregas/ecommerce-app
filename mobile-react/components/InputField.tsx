import { Colors } from "@/constants/Colors";
import React from "react";
import { StyleSheet, TextInput, TextInputProps, Text } from "react-native";

interface InputFieldProps extends TextInputProps {
  error?: string;
}

const InputField = ({ error, style, ...props }: InputFieldProps) => {
  return (
    <>
      <TextInput
        {...props}
        placeholderTextColor={Colors.gray}
        style={[
          styles.inputField,
          props.multiline && styles.multilineInput,
          style,
          error && styles.errorInput
        ]}
      />
      {error && <Text style={styles.errorText}>{error}</Text>}
    </>
  );
};

export default InputField;

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
    height: 48,
  },
  multilineInput: {
    height: undefined,
    minHeight: 60,
    textAlignVertical: 'top',
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
});
