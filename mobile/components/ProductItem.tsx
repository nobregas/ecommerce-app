import { Colors } from "@/constants/Colors";
import { ProductStrType, ProductType } from "@/types/type";
import { Ionicons } from "@expo/vector-icons";
import { Link } from "expo-router";
import { View, Text, StyleSheet, Dimensions, TouchableOpacity } from "react-native";
import Animated, { FadeInDown } from "react-native-reanimated";
import { Image } from "expo-image"
import React, { useEffect } from "react";
import productService, { SimpleProductObject } from "@/service/productService";
import HeartButton from "./HeartButton";

type Props = {
    item: SimpleProductObject;
    index: number
    productType: ProductStrType
}

const width = Dimensions.get("window").width - 40;

export default function ProductItem({ item, index, productType }: Props) {
    const [isFavorited, setIsFavorited] = React.useState(item.isFavorite);

    useEffect(() => {
        setIsFavorited(item.isFavorite)
    }, [isFavorited])

    return (
        <Link href={{
            pathname: "/product-details/[id]",
            params: { id: item.id, productType: productType }
        }} asChild>
            <TouchableOpacity>
                <Animated.View style={styles.container} entering={FadeInDown.delay(300 + index * 100).duration(500)}>
                    <Image
                        source={{ uri: item.image }}
                        style={styles.productImage}
                        placeholder={"https://via.placeholder.com/150"}
                    />
                    <HeartButton
                        productId={item.id}
                        initialIsFavorited={item.isFavorite}
                        size={20}
                        style={styles.bookmarkBtn}
                    />
                    <View style={styles.productInfo}>
                        <Text style={styles.price}>R${item.price}</Text>
                        <View style={styles.ratingWrapper}>
                            <Ionicons name="star" size={20} color={Colors.yellow} />
                            <Text style={styles.rating}>{item.averageRating}</Text>
                        </View>
                    </View>
                    <Text style={styles.title}>{item.title}</Text>
                </Animated.View>
            </TouchableOpacity>

        </Link>

    )
}

const styles = StyleSheet.create({
    container: {
        width: width / 2 - 10,

    },
    productImage: {
        width: "100%",
        height: 200,
        borderRadius: 15,
        marginBottom: 10
    },
    
    title: {
        fontSize: 14,
        fontWeight: "600",
        color: Colors.black,
        letterSpacing: 1.1
    },
    productInfo: {
        flexDirection: "row",
        justifyContent: "space-between",
        marginBottom: 8
    },
    price: {
        fontSize: 16,
        fontWeight: "700",
        color: Colors.primary
    },
    ratingWrapper: {
        flexDirection: "row",
        alignItems: "center",
        gap: 5
    },
    rating: {
        fontSize: 14,
        color: Colors.gray
    },
    bookmarkBtn: {
        position: "absolute",
        right: 20,
        top: 20,
        backgroundColor: Colors.transparentWhite,
        padding: 5,
        borderRadius: 30
    },
});