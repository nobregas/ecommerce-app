import { StyleSheet, Text, TouchableOpacity, View, KeyboardAvoidingView, Platform, ScrollView } from 'react-native'
import React from 'react'
import { Link, router, Stack } from 'expo-router'
import { Ionicons } from '@expo/vector-icons'
import { Colors } from '@/constants/Colors'
import InputField from '@/components/InputField'
import SocialLoginButtons from '@/components/SocialLoginButtons'

type Props = {}

const SignUpScreen = (props: Props) => {
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
          placeholderTextColor={Colors.gray}
          autoCapitalize="words"
        />
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
        <View style={styles.socialLoginWrapper}>
          <SocialLoginButtons emailHref={"/signin"} />
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
  }
});
