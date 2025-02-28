import { StyleSheet, Text, TextInput, TouchableOpacity, View } from 'react-native'
import React from 'react'
import { Link, router, Stack } from 'expo-router'
import { Ionicons } from '@expo/vector-icons'
import { Colors } from '@/constants/Colors'
import InputField from '@/components/InputField'
import SocialLoginButtons from '@/components/SocialLoginButtons'

type Props = {}

const SignUpScreen = (props: Props) => {
  return (
    <>
      <Stack.Screen options={{
        headerTitle: 'SignUp',
        headerTitleAlign: 'center',
        headerLeft: () => (
          <TouchableOpacity onPress={() => router.back()}>
            <Ionicons name='close' size={24} color={Colors.black} />
          </TouchableOpacity>
        )
      }} />
      <View style={styles.container}>
        <Text style={styles.title} >Create an account</Text>
        <InputField
          placeholder='Email Address'
          placeholderTextColor={Colors.gray}
          autoCapitalize='none'
          keyboardType='email-address'
        />
        <InputField
          placeholder='Password'
          placeholderTextColor={Colors.gray}
          secureTextEntry
        />
        <InputField
          placeholder='Confirm Password'
          placeholderTextColor={Colors.gray}
          secureTextEntry
        />

        <TouchableOpacity style={styles.btn}>
          <Text style={styles.btnTxt}>Create an Account</Text>
        </TouchableOpacity>

        <View style={styles.loginTxtWrapper}>
          <Text style={styles.loginTxt}>
            Already have an account? {" "}
            <Link href={"/signin"} asChild>
              <Text style={styles.loginTxtSpan}>SignIn</Text>
            </Link>
          </Text>
        </View>

        <View style={styles.divider} />

        <SocialLoginButtons emailHref={"/signin"} />
      </View>
    </>
  )
}

export default SignUpScreen

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