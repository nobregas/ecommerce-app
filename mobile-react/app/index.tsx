import { ImageBackground, StyleSheet, Text, View } from "react-native";
import React from "react";
import { Link, Stack } from "expo-router";
import { LinearGradient } from "expo-linear-gradient";
import { Colors } from "@/constants/Colors";
import Animated, { FadeInRight } from "react-native-reanimated";
import SocialLoginButtons from "@/components/SocialLoginButtons";

type Props = {};

const WelcomeScreen = (props: Props) => {
  return (
    <>
      <Stack.Screen options={{ headerShown: false }} />
      <ImageBackground
        source={require("@/assets/images/ecommerce-splash.jpg")}
        style={{ flex: 1 }}
        resizeMode="cover"
      >
        <View style={styles.container}>
          <LinearGradient
            colors={["transparent", Colors.semiTransparentWhite, Colors.solidWhite]}
            style={styles.background}
          >
            <View style={styles.wrapper}>
              <Animated.Text
                style={styles.title}
                entering={FadeInRight.delay(500).duration(300).springify()}
              >
                ShopX
              </Animated.Text>
              <Animated.Text
                style={styles.description}
                entering={FadeInRight.delay(500).duration(300).springify()}
              >
                One Stop Solution for All You Needs.
              </Animated.Text>

              <SocialLoginButtons emailHref={"/(auth)/signup"} />

              <View style={styles.loginTxtWrapper}>
                <Text style={styles.loginTxt}>
                  Already have an account? {" "}
                  <Link href={"/(auth)/signin"} asChild>
                    <Text style={styles.loginTxtSpan}>SignIn</Text>
                  </Link>
                </Text>
              </View>

            </View>
          </LinearGradient>
        </View>
      </ImageBackground>
    </>
  );
};

export default WelcomeScreen;

const styles = StyleSheet.create({
  container: {
    flex: 1,
    justifyContent: "center",
    alignItems: "center",
  },
  background: {
    flex: 1,
    position: "absolute",
    top: 0,
    bottom: 0,
    left: 0,
    right: 0,
    justifyContent: "flex-end",
  },
  wrapper: {
    paddingBottom: 50,
    paddingHorizontal: 20,
    alignItems: "center",
  },
  title: {
    fontSize: 22,
    color: Colors.primary,
    fontWeight: "700",
    letterSpacing: 2.4,
    marginBottom: 5,
  },
  description: {
    fontSize: 14,
    color: Colors.gray,
    letterSpacing: 1.2,
    lineHeight: 30,
    marginBottom: 20,
  },

  loginTxtWrapper: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "center",
    flexWrap: "wrap",
  },
  loginTxt: {
    margin: 30,
    fontSize: 14,
    color: Colors.black,
    lineHeight: 24,

  },
  loginTxtSpan: {
    color: Colors.primary,
    fontWeight: "600",

  },
});
