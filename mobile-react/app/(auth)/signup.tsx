import { StyleSheet, Text, TouchableOpacity, View, KeyboardAvoidingView, Platform, ScrollView } from 'react-native'
import React from 'react'
import { Link, router, Stack } from 'expo-router'
import { Ionicons } from '@expo/vector-icons'
import { Colors } from '@/constants/Colors'
import InputField from '@/components/InputField'
import SocialLoginButtons from '@/components/SocialLoginButtons'
import { isValidCPF, isValidEmail } from '@/utils/sharedFunctions'
import { useAuth } from '@/hooks/useAuth'

type Props = {}

const SignUpScreen = (props: Props) => {
  const [fullName, setFullName] = React.useState({ value: "", dirty: false });
  const [email, setEmail] = React.useState({ value: "", dirty: false });
  const [cpf, setCpf] = React.useState({ value: "", dirty: false });
  const [password, setPassword] = React.useState({ value: "", dirty: false });
  const [confirmPassword, setConfirmPassword] = React.useState({ value: "", dirty: false });

  const { handleAuth, loading, error } = useAuth()

  const validateFields = () => {
    let isValid = true

    if (!fullName.value.trim()) {
      setFullName(prev => ({ ...prev, dirty: true }))
      isValid = false
    }

    if (!isValidEmail(email.value)) {
      setEmail(prev => ({ ...prev, dirty: true }));
      isValid = false;
    }
    if (!isValidCPF(cpf.value)) {
      setCpf(prev => ({ ...prev, dirty: true }));
      isValid = false;
    }

    if (password.value.length < 3) {
      setPassword(prev => ({ ...prev, dirty: true }));
      isValid = false;
    }

    if (password.value !== confirmPassword.value) {
      setConfirmPassword(prev => ({ ...prev, dirty: true }));
      isValid = false;
    }

    return isValid
  }

  const getEmailError = () => {
    if (!email.dirty) return ""
    if (!email.value) return "Email is required"
    if (!isValidEmail(email.value)) return "Invalid email"
    return ""
  }

  const getCpfError = () => {
    if (!cpf.dirty) return '';
    if (!cpf.value) return 'CPF is required';
    if (!isValidCPF(cpf.value)) return 'Invalid CPF';
    return '';
  };

  const handleSubmit = async () => {
    if (!validateFields()) return

    await handleAuth("register", {
      fullName: fullName.value,
      email: email.value,
      cpf: cpf.value,
      password: password.value
    }, "/(auth)/signin");
  }

  return (
    <KeyboardAvoidingView
      style={{ flex: 1, paddingBottom: 60 }}
      behavior={Platform.OS === 'ios' ? 'padding' : 'height'}
    >
      <ScrollView
        contentContainerStyle={styles.container}
        keyboardShouldPersistTaps="handled"
      >
        <Stack.Screen options={{
          headerTitle: 'SignUp',
          headerTitleAlign: 'center',
          headerLeft: () => (
            <TouchableOpacity onPress={() => router.back()}>
              <Ionicons name='close' size={24} color={Colors.black} />
            </TouchableOpacity>
          )
        }} />
        <Text style={styles.title}>Create an account</Text>
        <InputField
          placeholder='Full Name'
          value={fullName.value}
          onChangeText={(text) => setFullName({ value: text, dirty: true })}
          error={fullName.dirty && !fullName.value ? 'Full name is required' : ''}
        />
        <InputField
          placeholder='Email Address'
          value={email.value}
          onChangeText={(text) => setEmail({ value: text, dirty: true })}
          error={getEmailError()}
        />
        <InputField
          placeholder='CPF'
          value={cpf.value}
          onChangeText={(text) => setCpf({ value: text, dirty: true })}
          error={getCpfError()}
        />
        <InputField
          placeholder='Password'
          value={password.value}
          onChangeText={(text) => setPassword({ value: text, dirty: true })}
          secureTextEntry={true}
          error={password.dirty && password.value.length < 3 ? 'Password must be at least 3 characters' : ''}
        />
        <InputField
          placeholder='Confirm Password'
          value={confirmPassword.value}
          onChangeText={(text) => setConfirmPassword({ value: text, dirty: true })}
          secureTextEntry={true}
          error={confirmPassword.dirty && password.value !== confirmPassword.value ? 'Passwords do not match' : ''}
        />

        {error && <Text style={styles.error}>{error}</Text>}
        <TouchableOpacity
          style={[styles.btn, loading && styles.disabledBtn]}
          onPress={handleSubmit}
          disabled={loading}
        >
          <Text style={styles.btnTxt}>
            {loading ? 'Loading...' : 'Create an Account'}
          </Text>
        </TouchableOpacity>

        <View style={styles.loginTxtWrapper}>
          <Text style={styles.loginTxt}>
            Already have an account? {" "}
            <Link href={"/(auth)/signin"} asChild>
              <Text style={styles.loginTxtSpan}>SignIn</Text>
            </Link>
          </Text>
        </View>

        <View style={styles.divider} />
        <View style={styles.socialLoginWrapper}>
          <SocialLoginButtons emailHref={"/(auth)/signin"} />
        </View>
      </ScrollView>
    </KeyboardAvoidingView>
  )
}

export default SignUpScreen

const styles = StyleSheet.create({
  container: {
    flexGrow: 1,
    justifyContent: 'flex-start',
    alignItems: 'center',
    padding: 20,
    backgroundColor: Colors.background,
    paddingBottom: 40,
  },
  title: {
    fontSize: 24,
    fontWeight: "600",
    letterSpacing: 1.2,
    color: Colors.black,
    marginBottom: 50,
  },
  btn: {
    backgroundColor: Colors.primary,
    paddingVertical: 14,
    paddingHorizontal: 18,
    alignSelf: 'stretch',
    alignItems: 'center',
    borderRadius: 5,
    marginBottom: 20,
  },
  btnTxt: {
    color: Colors.white,
    fontSize: 16,
    fontWeight: '600'
  },
  loginTxtWrapper: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "center",
    flexWrap: "wrap",
  },
  loginTxt: {
    marginBottom: 30,
    fontSize: 14,
    color: Colors.black,
    lineHeight: 24,
  },
  loginTxtSpan: {
    color: Colors.primary,
    fontWeight: "600",
  },
  divider: {
    borderTopColor: Colors.gray,
    borderTopWidth: StyleSheet.hairlineWidth,
    width: '30%',
    marginBottom: 30
  },
  socialLoginWrapper: {
    width: '100%',
    alignItems: 'center',
    marginTop: 20,
  },
  error: {
    color: 'red',
    fontSize: 14,
    marginBottom: 10,
  },
  disabledBtn: {
    opacity: 0.7,
  }
});
