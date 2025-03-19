import React, { memo } from "react";
import { View, TouchableOpacity, FlatList, Text, StyleSheet } from "react-native";
import ProductItem from "./ProductItem";
import { Colors } from "@/mobile/constants/Colors";
import { ProductType } from "@/mobile/types/type";

type Props = {
    products: ProductType[];
    flatlist?: boolean;
};

const productType = "regular";

const ProductList = ({ products, flatlist = true }: Props) => {
    const renderItem = ({ item, index }: { item: ProductType; index: number }) => (
        <ProductItem item={item} index={index} productType={productType} />
    );

    return (
        <View style={styles.container}>
            <View style={styles.titleWrapper}>
                <Text style={styles.title}>For you</Text>
                <TouchableOpacity>
                    <Text style={styles.titleBtn}>See All</Text>
                </TouchableOpacity>
            </View>
            {flatlist ? (
                <FlatList
                    data={products}
                    numColumns={2}
                    columnWrapperStyle={styles.columnWrapper}
                    keyExtractor={(item) => item.id.toString()}
                    renderItem={renderItem}
                />
            ) : (
                <View style={styles.itemsWrapper}>
                    {products.map((item, index) => (
                        <View key={item.id} style={styles.productWrapper}>
                            <ProductItem item={item} index={index} productType={productType} />
                        </View>
                    ))}
                </View>
            )}
        </View>
    );
};

const styles = StyleSheet.create({
    container: {
        marginHorizontal: 20,
    },
    titleWrapper: {
        flexDirection: "row",
        justifyContent: "space-between",
        marginBottom: 10,
    },
    title: {
        fontSize: 18,
        fontWeight: "600",
        letterSpacing: 0.6,
        color: Colors.black,
    },
    titleBtn: {
        fontSize: 14,
        fontWeight: "500",
        letterSpacing: 0.6,
        color: Colors.black,
    },
    itemsWrapper: {
        width: "100%",
        flexDirection: "row",
        flexWrap: "wrap",
        alignItems: "stretch",
    },
    productWrapper: {
        width: "50%",
        paddingLeft: 5,
        marginBottom: 20,
    },
    columnWrapper: {
        justifyContent: "space-between",
        marginBottom: 20,
    },
});

export default memo(ProductList);