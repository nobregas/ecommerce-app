import ImageSlider from "@/components/ImageSlider";
import { Colors } from "@/constants/Colors";
import { ProductStrType, ProductType } from "@/types/type";
import { Ionicons } from "@expo/vector-icons";
import { router, Stack, useLocalSearchParams } from "expo-router";
import { useEffect, useState } from "react";
import { ActivityIndicator, Text, View, StyleSheet, TouchableOpacity, ScrollView } from "react-native";
import { useHeaderHeight } from "@react-navigation/elements";
import Animated, { FadeInDown, SlideInDown } from "react-native-reanimated";
import { useCartStore } from "@/store/cardBadgeStore";
import React from "react";
import productService, { ProductDetailsType } from "@/service/productService";
import HeartButton from "@/components/HeartButton";

export default function ProductDetails() {
  const { id, productType: productType } = useLocalSearchParams()
  const headerHeight = useHeaderHeight();

  const productId = Array.isArray(id) ? id[0] : id
  const productTypeStr = Array.isArray(productType) ? productType[0] : productType

  const [product, setProduct] = useState<ProductDetailsType | null>(null)
  const [loading, setLoading] = useState<boolean>(true)
  const { addToCart } = useCartStore();

  useEffect(() => {
    fetchProductDetails()
  }, [])

  const fetchProductDetails = async () => {
    if (!productId || !productTypeStr) return;

    try {
      const data = await productService.getProductDetails(parseInt(productId));
      setProduct(data);
    } finally {
      setLoading(false);
    }
  };

  const handleAddToCart = async () => {
    try {
      if (!product) return;

      await addToCart(product.id);
      alert("Product added to the cart!");
    } catch (error) {
      console.log(error)
      alert("Error adding product to the cart ");
    }
  }

  const getProductImages = (): string[] => {
    const imagesUrls: string[] = product?.images?.map((image: any) => image.imageUrl) || [];
    return imagesUrls
  }

  if (loading) {
    return (
      <>
        <Stack.Screen options={{
          title: 'Product Details',
          headerTitleAlign: 'center',
          headerTransparent: true,
          headerLeft: () => (
            <TouchableOpacity onPress={() => router.back()}>
              <Ionicons name="arrow-back" size={24} color={Colors.black} />
            </TouchableOpacity>
          ),
          headerRight: () => (
            <TouchableOpacity onPress={() => router.push('/cart')}>
              <Ionicons name="cart-outline" size={24} color={Colors.black} />
            </TouchableOpacity>
          )
        }} />
        <View>
          <ActivityIndicator size={"large"} />
        </View>
      </>
    )
  }

  const numberOfRating = 100


  return (
    <>
      <Stack.Screen options={{
        title: 'Product Details',
        headerTitleAlign: 'center',
        headerTransparent: true,
        headerLeft: () => (
          <TouchableOpacity onPress={() => router.back()}>
            <Ionicons name="arrow-back" size={24} color={Colors.black} />
          </TouchableOpacity>
        ),
        headerRight: () => (
          <TouchableOpacity onPress={() => router.push('/cart')}>
            <Ionicons name="cart-outline" size={24} color={Colors.black} />
          </TouchableOpacity>
        )
      }} />
      <ScrollView style={{ marginTop: headerHeight, marginBottom: 90 }}>
        {product && (
          <Animated.View entering={FadeInDown.delay(300).duration(500)}>
            <ImageSlider images={getProductImages()} />
          </Animated.View>
        )}
        {product && (
          <View style={styles.container}>
            <Animated.View style={styles.productIcons} entering={FadeInDown.delay(500).duration(500)}>
              <View style={styles.productIcons}>
                <Ionicons name="star" size={18} color={Colors.yellow} />
                <Text style={styles.rating}>
                  {product.averageRating}
                </Text>
              </View>
              <HeartButton productId={product.id} initialIsFavorited={product.isFavorite} />
            </Animated.View >

            <Animated.Text
              style={styles.title}
              entering={FadeInDown.delay(700).duration(500)}
            >
              {product.title}
            </Animated.Text>

            <Animated.View style={styles.priceWrapper} entering={FadeInDown.delay(900).duration(500)}>
              <Text style={styles.priceTxt}>R${product.price.toFixed(2)}</Text>
              {product.discountPercentage > 0 && (
                <>
                  <View style={styles.discountWrapper}>
                    <Text style={styles.discount}>{product.discountPercentage}%
                    </Text>
                  </View>
                  <Text style={styles.oldPrice}>R${product.basePrice}</Text>
                </>
              )}
            </Animated.View>

            <Animated.Text
              style={styles.description}
              entering={FadeInDown.delay(1100).duration(500)}
            >
              {product.description}
            </Animated.Text>
          </View>
        )}

      </ScrollView>
      <Animated.View
        style={styles.buttonWrapper}
        entering={SlideInDown.delay(500).duration(500)}
      >
        <TouchableOpacity
          style={[
            styles.btn,
            { backgroundColor: Colors.white, borderColor: Colors.primary, borderWidth: 1 }
          ]}
          onPress={handleAddToCart}
        >
          <Ionicons name="cart-outline" size={22} color={Colors.primary} />
          <Text style={[styles.btnTxt, { color: Colors.primary }]}>Add to Cart</Text>
        </TouchableOpacity>
        <TouchableOpacity style={styles.btn}>
          <Text style={styles.btnTxt}>Buy Now</Text>
        </TouchableOpacity>
      </Animated.View>
    </>
  );
}

const styles = StyleSheet.create({
  container: {
    paddingHorizontal: 20,
  },
  productIcons: {
    flexDirection: "row",
    alignItems: "center",
    justifyContent: "space-between",
    marginBottom: 5
  },
  rating: {
    marginLeft: 5,
    fontSize: 14,
    fontWeight: "400",
    color: Colors.gray
  },
  title: {
    fontSize: 20,
    fontWeight: "400",
    color: Colors.black,
    letterSpacing: 0.6,
    lineHeight: 32,
  },
  priceWrapper: {
    flexDirection: "row",
    alignItems: "center",
    marginTop: 10,
    gap: 5
  },
  priceTxt: {
    fontSize: 18,
    fontWeight: "600",
    color: Colors.black
  },
  discountWrapper: {
    backgroundColor: Colors.extraLightGray,
    padding: 5,
    borderRadius: 5,
  },
  discount: {
    fontSize: 14,
    fontWeight: "400",
    color: Colors.primary
  },
  oldPrice: {
    fontSize: 16,
    fontWeight: "400",
    textDecorationLine: "line-through",
    color: Colors.gray
  },
  description: {
    marginTop: 10,
    fontSize: 16,
    color: Colors.black,
    letterSpacing: 0.6,
    lineHeight: 24
  },
  productVariationWrapper: {
    flexDirection: "row",
    marginTop: 20,
    flexWrap: "wrap",
  },
  productVariationType: {
    width: "50%",
    gap: 5,
    marginBottom: 10,
  },
  productVariationTitle: {
    fontSize: 16,
    fontWeight: "500",
    color: Colors.black
  },
  productVariationValueWrapper: {
    flexDirection: "row",
    alignItems: "center",
    gap: 5,
    flexWrap: "wrap"
  },
  productVariationColorValue: {
    width: 30,
    height: 30,
    borderRadius: 15,
    backgroundColor: Colors.extraLightGray
  },
  productVariationSizeValue: {
    width: 50,
    height: 30,
    borderRadius: 5,
    backgroundColor: Colors.extraLightGray,
    justifyContent: "center",
    alignItems: "center",
    borderColor: Colors.lightGray,
    borderWidth: 1
  },
  productVariationSizeValueTxt: {
    fontSize: 12,
    fontWeight: "500",
    color: Colors.black
  },
  buttonWrapper: {
    position: "absolute",
    height: 90,
    padding: 20,
    bottom: 0,
    width: "100%",
    backgroundColor: Colors.white,
    flexDirection: "row",
    gap: 10,
  },
  btn: {
    flex: 1,
    backgroundColor: Colors.primary,
    height: 40,
    justifyContent: "center",
    alignItems: "center",
    borderRadius: 5,
    gap: 5,
    flexDirection: "row",
    elevation: 5,
    shadowColor: Colors.black,
    shadowOffset: {
      width: 0,
      height: 2,
    },
    shadowOpacity: 0.25,
    shadowRadius: 3.84,
  },
  btnTxt: {
    fontSize: 16,
    fontWeight: "500",
    color: Colors.white
  },
});