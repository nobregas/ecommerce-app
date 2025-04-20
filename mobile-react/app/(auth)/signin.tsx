import { StyleSheet, Text, TouchableOpacity, View } from 'react-native'
import React from 'react'
import { Link, router, Stack } from 'expo-router'
import InputField from '@/components/InputField'
import SocialLoginButtons from '@/components/SocialLoginButtons'
import { Ionicons } from '@expo/vector-icons'
import { Colors } from '@/constants/Colors'
import { isValidEmail } from '@/utils/sharedFunctions'
import { useAuth } from '@/hooks/useAuth'

type Props = {}

const SignInScreen = (props: Props) => {
  const [email, setEmail] = React.useState({ value: '', dirty: false });
  const [password, setPassword] = React.useState({ value: '', dirty: false });
  const { handleAuth, loading, error } = useAuth();

  const getEmailError = () => {
    if (!email.dirty) return '';
    if (!email.value) return 'Email is required';
    if (!isValidEmail(email.value)) return 'Invalid email format';
    return '';
  };

  const validateFields = () => {
      let isValid = true
  
      if (!isValidEmail(email.value)) {
        setEmail(prev => ({ ...prev, dirty: true }));
        isValid = false;
      }
      if (password.value.length < 3) {
        setPassword(prev => ({ ...prev, dirty: true }));
        isValid = false;
      }
  
      return isValid
    }

  const handleLogin = async () => {
    if (!validateFields()) return;
    
    await handleAuth('login', {
      email: email.value,
      password: password.value
    });
  };

  return (
    <>
      <Stack.Screen options={{
        headerTitle: 'SignIn',
        headerTitleAlign: 'center',
        headerLeft: () => (
          <TouchableOpacity onPress={() => router.back()}>
            <Ionicons name='close' size={24} color={Colors.black} />
          </TouchableOpacity>
        )
      }} />
      <View style={styles.container}>
        <Text style={styles.title} >Login to Your Account</Text>
        <InputField
          placeholder='Email Address'
          value={email.value}
          onChangeText={(text) => setEmail({ value: text, dirty: true })}
          error={getEmailError()}
        />
        <InputField
          placeholder='Password'
          value={password.value}
          onChangeText={(text) => setPassword({ value: text, dirty: true })}
          secureTextEntry={true}
          error={password.dirty && !password.value ? 'Password is required' : ''}
        />

        <TouchableOpacity style={styles.btn} onPress={handleLogin}>
          <Text style={styles.btnTxt}>Login</Text>
        </TouchableOpacity>

        <View style={styles.loginTxtWrapper}>
          <Text style={styles.loginTxt}>
            Don't have an account? {" "}
            <Link href={"/signup"} asChild>
              <Text style={styles.loginTxtSpan}>SignUp</Text>
            </Link>
          </Text>
        </View>

        <View style={styles.divider} />

        <SocialLoginButtons emailHref={"/signin"} />
      </View>
    </>


  )
}

export default SignInScreen

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    padding: 20,
    backgroundColor: Colors.background
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
  }
})